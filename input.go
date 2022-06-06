package main

import (
	"log"
)

const (
	excelSource ioCandidate = iota
	csvSource
	nonExistingSource
)

type ioSource struct {
	source   ioCandidate
	fileName string
	context  map[string]string
}

type ioCandidate uint

func (ic ioCandidate) New(file string, ctx map[string]string) ioSource {
	return ioSource{ic, file, ctx}
}

func ReadSource(ctx globalContext) Rows {
	result := make(Rows, 0)

	switch ctx.input.source {
	case excelSource:
		excelMaps := excelExtractor(ctx.input)
		result = mapsToRows(excelMaps, ctx)
		//for _, v := range excelMaps {
		//	row := excelConvertRow(v, ctx.input.context)
		//	result = append(result, row)
		//}
		return result
	}
	log.Fatal("input source not implemented")

	return result
}
