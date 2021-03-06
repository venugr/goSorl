package main

import (
	"errors"
	"fmt"
	"strings"
)

func sorlLoadGlobalVars(homePath string, svMap *SorlMap) error {

	sorlDefaultVarFile := homePath + PathSep + ".sorl" + PathSep + "vars.sorl"
	return readVarsFile(sorlDefaultVarFile, svMap)

}

func readVarsFile(fileName string, svMap *SorlMap) error {

	fileOrDir, err := chkFileOrDir(fileName)

	if err != nil {
		return err
	}

	if fileOrDir {
		fileName = fileName + PathSep + "vars.sorl"
	}

	if !chkDir(fileName) {
		return errors.New("\nfile is not found")
	}

	fmt.Println("\ninfo: reading vars file:", fileName)
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
