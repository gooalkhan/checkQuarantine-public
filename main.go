package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

//todo list
//1. BL번호 등 필수항목이 빠진 map은 row로 변환하지 말기
//2. 통상적이지 않은 에러(복수 동물검역 합격, 식품검역수가 아이템수랑 맞지않는 경우 등)를 개별 아이템 에러에도 넣고, 화면에도 출력
//3. toml 파서 만들기
//3.1 - 글로벌 변수 설정값 중에서는 재컴파일로 바꿔야 할 변수가 있고, 설정파일로 바꿔야 하는 변수가 있음
//4. GUI(웹팩?) 구현

func main() {
	defer func() {
		errorMsg := recover()
		if errorMsg != nil {
			log.Println(errorMsg)
		}
		fmt.Print("Press 'Enter' to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		err := DeleteTempFile()
		if err != nil {
			log.Printf("%s: cannot remove temp files", err)
		}
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()
	fmt.Println("Starting Quarantine Result Checker...")
	inputSource := FileProber(os.Args)
	fmt.Printf("the file name is %s\n", inputSource.fileName)

	outFileName := TempFileNameGenerator(excelSource)
	outputSource := excelSource.New(outFileName, excelContext)
	con := NewContext(unipassToken, inputSource, outputSource, 3)

	rows := ReadSource(con)

	cntrs := rows.toContainers(con)
	cntrs.GetUnipass(con)
	cntrs.checkAllConditions()

	cntrs.WriteSource(con)

	// convertedRows := cntrs.toRows()
	// if inputSource.source == excelSource && outputSource.source == excelSource {
	// 	err := convertedRows.excelAddToOriginWriter(inputSource, outputSource.fileName)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	ole.CoInitialize(0)
	// 	unknown, _ := oleutil.CreateObject("Excel.Application")
	// 	excel, _ := unknown.QueryInterface(ole.IID_IDispatch)
	// 	oleutil.PutProperty(excel, "Visible", true)

	// 	workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()

	// 	cwd, err := os.Getwd()

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	readExample(cwd+"\\"+outFileName, excel, workbooks)
	// 	workbooks.Release()
	// 	//oleutil.CallMethod(excel, "Quit")
	// 	excel.Release()
	// 	ole.CoUninitialize()

	// } else {
	// 	panic(notImplementError)
	// }
}
