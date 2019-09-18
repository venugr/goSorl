package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func sorlLoadGlobalVars(homePath string, svMap *SorlMap) error {

	sorlDefaultVarFile := homePath + PathSep + ".sorl" + PathSep + "vars.sorl"
	return readVarsFile(sorlDefaultVarFile, svMap)

}

func sorlLoadFileVars(varFileName string, svMap *SorlMap) error {

	return readVarsFile(varFileName, svMap)
}

func sorlArgsVars(svMap *SorlMap) error {

	cmdCnt := 1
	(*svMap)["_cmd.arg.order"] = ""

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-var=") {
			arg = strings.TrimPrefix(arg, "-var=")
			arg = strings.TrimPrefix(arg, "\"")
			arg = strings.TrimSuffix(arg, "\"")
			argList := strings.Split(arg, "=")
			arg = strings.TrimSuffix(arg, argList[0])
			(*svMap)[argList[0]] = arg
		}

		if strings.HasPrefix(arg, "-cmd=") {
			arg = strings.TrimPrefix(arg, "-cmd=")
			arg = strings.TrimPrefix(arg, "\"")
			arg = strings.TrimSuffix(arg, "\"")
			(*svMap)["_cmd.arg."+strconv.Itoa(cmdCnt)] = arg
			(*svMap)["_cmd.arg.order"] += "_cmd.arg." + strconv.Itoa(cmdCnt) + ","
			cmdCnt++
		}

	}

	return nil
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
