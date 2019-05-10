package main

import (
	"fmt"
	"strings"
)

func main() {

	fmt.Println()

	cliArgsMap := getCliArgs()
	fmt.Println(cliArgsMap)

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	fmt.Println(envMap)

	homePath := envMap["HOME"]
	userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	globalOrchFilePath := strings.TrimSpace(cliArgsMap["orchfile"])

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

	sorlStart(globalOrchFilePath, scProp, hostList)

	fmt.Println()
	fmt.Println()
}
