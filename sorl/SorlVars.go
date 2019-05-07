package main

import (
	"errors"
	"fmt"
	"strings"
)

func readVarsFile(fileName string, svMap *SorlMap) error {

	fileOrDir, err := chkFileOrDir(fileName)

	if err != nil {
		return err
	}

	if fileOrDir {
		fileName = fileName + "/" + "vars.sorl"
	}

	if !chkDir(fileName) {
		return errors.New("\nfile is not found")
	}

	fmt.Println("info: reading vars file:", fileName)
	lines, _ := ReadFile(fileName)
	idx := -1
	varKey := ""
	varVal := ""

	for _, line := range lines {

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		if strings.Contains(line, "=") {
			idx = strings.Index(line, "=")
			varKey = line[:idx]
			varVal = line[idx+1:]
			(*svMap)[varKey] = varVal
		}
	}

	return nil
}
