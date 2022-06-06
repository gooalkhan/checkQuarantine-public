package main

import "checkQuarantine/env"

const (
	unipassToken string = env.UnipassToken
	unipassURL   string = "https://unipass.customs.go.kr:38010/ext/rest/cargCsclPrgsInfoQry/retrieveCargCsclPrgsInfo?"
)

//var dirPath string = ""

type globalContext struct {
	token     string
	url       string
	semaphore int64
	input     ioSource
	output    ioSource
}

var excelContext = map[string]string{
	"headerRowNumber": headerRowNumber,
	"sheetName":       sheetName,
	"order":           excelOrder,
	"bl":              excelBl,
	"eta":             excelEta,
	"etd":             excelEtd,
	"code":            excelCode,
	"name":            excelName,
	"qty":             excelQty,
}

func NewContext(token string, input ioSource, output ioSource, sem int64) globalContext {
	result := globalContext{}

	result.input = input
	result.output = output
	result.token = token
	result.url = unipassURL
	result.semaphore = sem

	return result
}
