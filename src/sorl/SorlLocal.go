package main

import (
	"fmt"
	"strings"
)

func sorlLocal(parallelOk, orchFile string, scProp SorlConfigProperty,
	hostsList []string, cliArgsMap map[string]string,
	svMap map[string]string) {

	allProp := Property{}

	for fKey, fVal := range svMap {
		allProp[fKey] = fVal
	}

	cmdStr := ""

	for _, lVal := range strings.Split(allProp["_cmd.arg.order"], ",") {

		cmdStr += allProp[lVal] + "\n"

	}

	cmdStr = strings.TrimSpace(cmdStr)

	//fmt.Println("WaitPrompt: " + waitPrompt)

	commands := []string{
		cmdStr,
		"",
	}

	fmt.Println(commands)

}
