package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {

	testProp := Property{
		"a":    "111",
		"b":    "222",
		"ab":   "1122",
		"test": "Soooful",
	}

	if false {
		//oLine := "abcd{hello}{world} {a{b{c{ d  }}}}this is a prop replace{name}{ lname    	}{doit{howtodo}}"
		oLine := "{test},one={a} and two={b}, chk the prop replace{ab}"
		mLine, err1 := replaceProp(oLine, testProp)
		checkError(err1)

		fmt.Printf("\n%s\n%s", oLine, mLine)
	}
	//.Println("Hello Pavana..Welcome to Go!")
	//return

	cliArgsMap := getCliArgs()
	//fmt.Println(cliArgsMap)
	printVersion()
	versionOk := strings.TrimSpace(cliArgsMap["version"])

	if versionOk == "true" {
		return
	}

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	//printMap("ENVIRONMENT", map[string]string(envMap))
	//logit("\n")

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
		logit(fmt.Sprintf("\ninfo: %v", err))
	}
	printMap("Global Vars", svMap)

	hostList, err := sorlProcessCliArgs(scProp, cliArgsMap)

	if err != nil {
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		return
	}
	PrintList("All the selected hosts", hostList)
	//os.Exit(1)
	sorlStart(parallelOk, globalOrchFilePath, scProp, hostList, cliArgsMap, svMap)

	fmt.Println()
	fmt.Println()
}

func printVersion() {
	fmt.Println()
	fmt.Println("SORL: Solution ORchestration Language, the .(dot) scripting")
	fmt.Println("Version: 0.1-beta-build-1.0")
	time.Sleep(3 * time.Second)
	fmt.Println()
}
