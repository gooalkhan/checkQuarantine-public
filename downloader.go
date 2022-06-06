package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/sync/semaphore"
)

func (cntr *aContainer) Get(ctx globalContext) {
	urlMbl := ctx.url + fmt.Sprintf("crkyCn=%s", ctx.token) + fmt.Sprintf("&mblNo=%s", cntr.bl) + fmt.Sprintf("&blYy=%d", cntr.year)
	urlHbl := ctx.url + fmt.Sprintf("crkyCn=%s", ctx.token) + fmt.Sprintf("&hblNo=%s", cntr.bl) + fmt.Sprintf("&blYy=%d", cntr.year)
	urls := []string{urlMbl, urlHbl}
	isDone := false

	for i := 0; !isDone && i <= 1; i++ {
		var temp retrievedXML
		resp, err := http.Get(urls[i])
		if err != nil {
			log.Fatal(err.Error())
		}
		defer func(resp *http.Response) {
			err := resp.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp)

		rawByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		temp = retrievedXML{}
		err = xml.Unmarshal(rawByte, &temp)
		if err != nil {
			log.Fatal(err)
		}
		cntr.retrievedXML = temp
		if temp.TCnt >= 1 {
			isDone = true
			break
		}

	}

}

func (cntrs Containers) GetUnipass(ctx globalContext) chan *aContainer {
	var w sync.WaitGroup
	c := make(chan *aContainer, len(cntrs))
	s := semaphore.NewWeighted(ctx.semaphore)

	for i := 0; i < len(cntrs); i++ {
		w.Add(1)
		err := s.Acquire(context.Background(), 1)
		if err != nil {
			log.Fatalln(err)
		}
		go func(cntr *aContainer, ctx globalContext, j int) {
			fmt.Printf("start %d\n", j)
			cntr.Get(ctx)
			fmt.Printf("got %d\n", j)
			c <- cntr
			fmt.Printf("%d sent to channel\n", j)
			s.Release(1)
			fmt.Printf("%d semaphore released\n", j)
			w.Done()
			fmt.Printf("%d waitgroup done\n", j)
		}(&cntrs[i], ctx, i)
	}
	w.Wait()
	close(c)

	return c
}
