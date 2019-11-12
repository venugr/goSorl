package main

import (
	"fmt"
	"strings"
)

func sorlRunScript(scriptName string, scProp SorlConfigProperty, cliArgsMap map[string]string) error {

	fileInfo, err := ReadFile(scriptName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	ss := &SorlSSH{}
	alp := Property{}
	alp["sr:debug"] = cliArgsMap["debug"]
	alp["sr:info"] = cliArgsMap["info"]
	ss.sorlOrchestration(strings.Join(fileInfo, "\n"), &alp)

	return nil

}
