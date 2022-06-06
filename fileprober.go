package main

import (
	"regexp"
)

type myErrorType int

const (
	notFileError myErrorType = iota
	manyArgsError
	noFileError
	notImplementError
)

func (err myErrorType) Error() string {
	switch err {
	case notFileError:
		return "이 파일의 확장자는 지원하지 않습니다.\n"
	case manyArgsError:
		return "한개의 파일만 입력해주세요\n"
	case noFileError:
		return "파일을 입력해주세요\n"
	case notImplementError:
		return "기능이 구현되지 않았거나 알수없는 에러입니다\n"
	default:
		return "error case not implemented\n"
	}
}

func errorBool(e error) bool {
	if e != nil {
		return true
	}
	return false
}

func checkFileType(s string) (ioCandidate, error) {
	re := regexp.MustCompile("\\.(\\w{3,4})$")
	exten := re.FindStringSubmatch(s)
	if len(exten) != 0 {
		switch exten[len(exten)-1] {
		case "xls", "xlsx":
			return excelSource, nil
		case "csv":
			return csvSource, nil
		default:
		}
	}
	return nonExistingSource, notFileError
}

func FileProber(args []string) ioSource {
	if len(args) != 2 {
		if len(args) > 2 {
			panic(manyArgsError)
		}
		panic(noFileError)
	}

	ioca, err := checkFileType(args[1])

	if err != nil {
		panic(err)
	}

	switch ioca {
	case excelSource:
		return ioca.New(args[1], excelContext)
	case csvSource:
		panic(notImplementError)
	}
	panic(notImplementError)
}
