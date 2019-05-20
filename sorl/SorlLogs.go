package main

import (
	"bytes"
	"fmt"
	"log"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "Logger: ", log.Lshortfile)
)

func logit(logStr string) {

	//logger.Print(logStr)
	fmt.Print(logStr)

}
