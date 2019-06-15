package main

import (
	"errors"
	"fmt"
)

func sorlStartOrchestration(orchFile, lHost string, varsPerHostMap SorlMap, scProp SorlConfigProperty) {

	fmt.Printf("\n\ninfo: start ochestration for host: %s", lHost)
	fmt.Printf("\n\tOrchestration File: %s", orchFile)

	lines, err := readOrchFile(orchFile)

	if err != nil {
		fmt.Printf("\nerror: file/path '%s' not found", orchFile)
	}

	PrintList("File: "+orchFile, lines)

}

func readOrchFile(fileName string) ([]string, error) {

	fileOrDir, err := chkFileOrDir(fileName)

	if err != nil {
		return nil, err
	}

	if fileOrDir {
		fileName = fileName + PathSep + "main.sorl"
	}

	if !chkDir(fileName) {
		return nil, errors.New("\nfile is not found")
	}

	fmt.Println("\ninfo: reading orch file:", fileName)

	lines, _ := ReadFile(fileName)

	return lines, nil

}
