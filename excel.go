package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	headerRowNumber string = "6"
	sheetName       string = "Sheet1"
	readFileName    string = "text.xlsx"
	writeFileName   string = "checked.xlsx"

	excelOrder string = "Order.No"
	excelBl    string = "B/L.No"
	excelEta   string = "E.T.A"
	excelEtd   string = "E.T.D"
	excelCode  string = "품목코드"
	excelName  string = "품명"
	excelQty   string = "Q`ty"

	excelQuarantine string = "[F|C]"
	excelFood       string = ".*"
	excelWarehouse  string = ".*"
	excelArrive     string = ".*"
)

func excelExtractor(is ioSource) []map[string]string {
	result := make([]map[string]string, 0)
	f, err := excelize.OpenFile(is.fileName)
	if err != nil {
		log.Fatalln(err)
	}

	excelHeaderRowNum, err := strconv.Atoi(is.context["headerRowNumber"])

	rows := f.GetRows(is.context["sheetName"])
	header, rows := rows[excelHeaderRowNum], rows[excelHeaderRowNum+1:]

	for k, row := range rows {
		var temp map[string]string
		temp = make(map[string]string)
		for i, v := range row {
			if _, ok := temp[header[i]]; ok {
				log.Fatal("header should be unique!")
			}
			temp[header[i]] = v
		}
		temp["Row"] = strconv.Itoa(k)
		result = append(result, temp)
	}
	return result
}

func (rows Rows) excelAddToOriginWriter(in ioSource, filename string) error {
	//initializing
	f, err := excelize.OpenFile(in.fileName)
	if err != nil {
		panic(err)
	}
	startRow, err := strconv.Atoi(headerRowNumber)
	startRow++
	startRowString := strconv.Itoa(startRow)
	if err != nil {
		log.Print(err)
		return err
	}
	sheetname := in.context["sheetName"]
	allRows := f.GetRows(sheetname)
	headerCol := allRows[startRow]

	//write header
	startCol := len(headerCol)
	var colAlpha string
	colAlpha = excelize.ToAlphaString(startCol)
	f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("필수조건"))
	colAlpha = excelize.ToAlphaString(startCol + 1)
	f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("필수조건완료일"))
	colAlpha = excelize.ToAlphaString(startCol + 2)
	f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("필수조건실패"))

	startCol += 3

	for _, v := range containerConditions {
		context, err := v.Context(in.source)
		if err != nil {
			log.Print(err)
			return err
		}
		colAlpha = excelize.ToAlphaString(startCol)
		f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("%s", context["name"]))
		startCol++
		colAlpha = excelize.ToAlphaString(startCol)
		f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("Date"))
		startCol++
	}

	for _, v := range itemConditions {
		context, err := v.Context(in.source)
		if err != nil {
			log.Print(err)
			return err
		}
		colAlpha = excelize.ToAlphaString(startCol)
		f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("%s", context["name"]))
		startCol++
		colAlpha = excelize.ToAlphaString(startCol)
		f.SetCellStr(sheetname, colAlpha+startRowString, fmt.Sprintf("Date"))
		startCol++
	}
	//write body
	for _, v := range rows {
		v.excelAddToOriginRowWriter(f, sheetname, len(headerCol))
	}

	err = f.SaveAs(filename)
	if err != nil {
		panic(err)
	}

	//when not error
	return nil
}

func (row Row) excelAddToOriginRowWriter(f *excelize.File, sheet string, startCol int) {
	var col string
	var rowNum string
	startRow, err := strconv.Atoi(headerRowNumber)
	if err != nil {
		log.Fatal(err)
	}
	startRow = startRow + 2

	//write must conditions
	mustbool, musttime, mustcond := row.mustConditionReducer(excelSource)
	rowNum = strconv.Itoa(row.rowNum + startRow)
	col = excelize.ToAlphaString(startCol)
	f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%t", mustbool))
	col = excelize.ToAlphaString(startCol + 1)
	f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%s", musttime.Format("2006-01-02")))
	col = excelize.ToAlphaString(startCol + 2)
	f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%v", mustcond))

	startCol = startCol + 3

	//write container conditions
	for _, k := range containerConditions {
		for _, v := range row.contConditions {
			rowNum = strconv.Itoa(row.rowNum + startRow)
			col = excelize.ToAlphaString(startCol + indexList(v.condName, containerConditions)*2)
			if k == v.condName {
				f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%t", v.isCondApproved))
				col = excelize.ToAlphaString(startCol + indexList(v.condName, containerConditions)*2 + 1)
				f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%s", v.dateCondApproved.Format("2006-01-02")))
			}
		}
	}

	//write item conditions
	for _, k := range itemConditions {
		for _, v := range row.itemConditions {
			if k == v.condName {
				col = excelize.ToAlphaString(startCol + len(containerConditions)*2 + indexList(v.condName, itemConditions)*2)
				f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%t", v.isCondApproved))
				col = excelize.ToAlphaString(startCol + len(containerConditions)*2 + indexList(v.condName, itemConditions)*2 + 1)
				f.SetCellStr(sheet, col+rowNum, fmt.Sprintf("%s", v.dateCondApproved.Format("2006-01-02")))
			}
		}
	}
}
