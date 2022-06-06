package main

import (
	"log"
	"os"
)

type logLevel int

var myLogger *log.Logger

const (
	CRITICAL logLevel = iota
	WARNING
	INFO
	DEBUG
)

func init() {
	myLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
