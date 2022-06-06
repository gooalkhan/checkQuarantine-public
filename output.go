package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func TempFileNameGenerator(io ioCandidate) string {
	now := time.Now()
	var filename string
	nowTimeStamp := now.Format("20060102150405")

	switch io {
	case excelSource:
		filename = "임시파일_" + nowTimeStamp + ".xlsx"
	case csvSource:
		filename = "임시파일_" + nowTimeStamp + ".csv"
	}
	return filename
}

func DeleteTempFile() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		return err
	}
	re := regexp.MustCompile("임시파일_\\d{14}\\.\\w{3,4}")
	for _, fi := range files {
		if re.MatchString(fi.Name()) {
			filename := fi.Name()
			os.Remove(filename)
			log.Printf("%s is removed\n", filename)
		}
	}
	return nil
}

func (cntrs Containers) WriteSource(ctx globalContext) error {
	rows := cntrs.toRows()

	switch ctx.output.source {
	case excelSource:
		switch ctx.input.source {
		case excelSource:
			err := rows.excelAddToOriginWriter(ctx.input, ctx.output.fileName)
			if err != nil {
				panic(err)
			}

			ole.CoInitialize(0)
			unknown, _ := oleutil.CreateObject("Excel.Application")
			excel, _ := unknown.QueryInterface(ole.IID_IDispatch)
			oleutil.PutProperty(excel, "Visible", true)

			workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()

			cwd, err := os.Getwd()

			if err != nil {
				panic(err)
			}

			readExample(cwd+"\\"+ctx.output.fileName, excel, workbooks)
			workbooks.Release()
			//oleutil.CallMethod(excel, "Quit")
			excel.Release()
			ole.CoUninitialize()
		default:
		}
	default:
	}
	return notImplementError
}
