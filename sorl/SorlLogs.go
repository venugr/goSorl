package main

import (
	"bytes"
	"log"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "Logger: ", log.Lshortfile)
)

func logit(logStr ...string) {

	logger.Print("Hello, Log file")
	logger.Print(logStr)

}
