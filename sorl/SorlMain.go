package main

import (
	"fmt"
	"strings"
)

func main() {

	testProp := Property{
		"a":    "111",
		"b":    "222",
		"ab":   "1122",
		"test": "Soooful",
	}

	//oLine := "abcd{hello}{world} {a{b{c{ d  }}}}this is a prop replace{name}{ lname    	}{doit{howtodo}}"
	oLine := "{test},one={a} and two={b}, chk the prop replace{ab}"
	mLine, err1 := replaceProp(oLine, testProp)
	checkError(err1)

	fmt.Printf("\n%s\n%s", oLine, mLine)

	fmt.Println()
	//return

	cliArgsMap := getCliArgs()
	fmt.Println(cliArgsMap)

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	fmt.Println(envMap)

	homePath := envMap["HOME"]
	userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	globalOrchFilePath := strings.TrimSpace(cliArgsMap["orchfile"])
	parallelOk := strings.TrimSpace(cliArgsMap["parallel"])

	scProp := SorlConfigProperty{}
	svMap := SorlMap{}

	sorlLoadConfigFiles(&scProp, homePath, userConfigFilePath)
	scProp.printConfig()

	err := sorlLoadGlobalVars(homePath, &svMap)

	if err != nil {
		fmt.Printf("\ninfo: %v", err)
	}
	printMap("Global Vars", svMap)

	hostList := sorlProcessCliArgs(scProp, cliArgsMap)
	PrintList("All the selected hosts", hostList)

	sorlStart(parallelOk, globalOrchFilePath, scProp, hostList, cliArgsMap, svMap)

	fmt.Println()
	fmt.Println()
}
