package main

import (
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func readExample(fileName string, excel, workbooks *ole.IDispatch) {
	workbook, err := oleutil.CallMethod(workbooks, "Open", fileName)

	if err != nil {
		panic(err)
	}
	defer workbook.ToIDispatch().Release()

	sheets := oleutil.MustGetProperty(excel, "Sheets").ToIDispatch()
	//sheetCount := (int)(oleutil.MustGetProperty(sheets, "Count").Val)
	//fmt.Println("sheet count=", sheetCount)
	sheets.Release()

	worksheet := oleutil.MustGetProperty(workbook.ToIDispatch(), "Worksheets", 1).ToIDispatch()
	defer worksheet.Release()
	// for row := 1; row <= 2; row++ {
	// 	for col := 1; col <= 5; col++ {
	// 		cell := oleutil.MustGetProperty(worksheet, "Cells", row, col).ToIDispatch()
	// 		val, err := oleutil.GetProperty(cell, "Value")
	// 		if err != nil {
	// 			break
	// 		}
	// 		fmt.Printf("(%d,%d)=%+v toString=%s\n", col, row, val.Value(), val.ToString())
	// 		cell.Release()
	// 	}
	// }
}
