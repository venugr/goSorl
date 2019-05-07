package main

import (
	"fmt"
	"os"
	"strings"
)

type SorlMap map[string]string

func main() {

	fmt.Println()

	cliArgsMap := getCliArgs()
	fmt.Println(cliArgsMap)

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	fmt.Println(envMap)

	homePath := envMap["HOME"]
	userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])

	scProp := SorlConfigProperty{}
	svMap := SorlMap{}

	sorlLoadConfigFiles(&scProp, homePath, userConfigFilePath)
	scProp.printConfig()

	err := sorlLoadGlobalVars(homePath, &svMap)

	if err != nil {
		fmt.Printf("\ninfo: %v", err)
	}
	printMap("Global Vars", svMap)

	hostList := sorlProcessArgs(scProp, cliArgsMap)
	PrintList("All the selected hosts", hostList)

	sorlStart(scProp, hostList)

	fmt.Println()
}

func sorlLoadGlobalVars(homePath string, svMap *SorlMap) error {

	sorlDefaultVarFile := homePath + "/" + ".sorl" + "/" + "vars.sorl"
	return readVarsFile(sorlDefaultVarFile, svMap)

}

func sorlProcessArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) []string {

	grpCliOk := false
	hostCliOk := false

	hostCli := strings.TrimSpace(cliArgsMap["host"])
	grpCli := strings.TrimSpace(cliArgsMap["group"])
	maxGoRout := cliArgsMap["max"]

	if hostCli != "" && grpCli != "" {
		fmt.Println("\nError: Both 'host' and 'group' can not be present.")
		fmt.Println()
		os.Exit(1)
	}

	if hostCli == "" && grpCli == "" {
		fmt.Println("\nError: One of 'host' or 'group' must be present.")
		fmt.Println()
		os.Exit(1)
	}

	hostGrpCli := hostCli
	hostCliOk = true
	selType := "host"

	if strings.TrimSpace(grpCli) != "" {
		hostGrpCli = grpCli
		selType = "group"
		grpCliOk = true
		hostCliOk = false
	}

	fmt.Println("Host/Group:", hostGrpCli)
	fmt.Println(grpCliOk, hostCliOk, selType, maxGoRout)

	hostList, _ := getHostList(selType, hostGrpCli, scProp)

	return hostList

}

func sorlLoadConfigFiles(scProp *SorlConfigProperty, homePath string, userConfigFilePath string) {

	sorlDefaultConfigFile := homePath + "/" + ".sorl" + "/" + "config.sorl"
	readConfigFilePath(scProp, sorlDefaultConfigFile)

	//userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	if userConfigFilePath != "" {
		if err := readConfigFilePath(scProp, userConfigFilePath); err != nil {
			fmt.Printf("Unable to read proceed, %v", err)
			os.Exit(1)
		}
	}

}

func getHostList(selType, hostGrpCli string, scProp SorlConfigProperty) ([]string, error) {
	hostList := []string{}

	if strings.EqualFold(selType, "host") && hostGrpCli != "all" {
		return strings.Split(hostGrpCli, ","), nil
	}

	if strings.EqualFold(selType, "host") && hostGrpCli == "all" {
		allHostNames := scProp["all.hosts"]
		keys := make([]string, len(allHostNames))
		i := 0
		for k := range allHostNames {
			keys[i] = k
			i++
		}

		return keys, nil
	}

	if strings.EqualFold(selType, "group") && hostGrpCli != "all" {
		mapKV := scProp["group.hosts"]
		hostsStr := ""

		for _, grpVal := range strings.Split(hostGrpCli, ",") {

			if mVal, ok := mapKV[grpVal]; ok {
				hostsStr = strings.TrimRight(hostsStr, ",")
				hostsStr = hostsStr + "," + mVal

			} else {
				return nil, fmt.Errorf("Error: group id '%s' is not found", hostGrpCli)
			}

		}
		hostsStr = strings.Trim(hostsStr, ",")

		return strings.Split(hostsStr, ","), nil

	}

	return hostList, nil
}

func readConfigFilePath(scProp *SorlConfigProperty, configFilePath string) error {
	/*
		if configFilePath == "" {
			return errors.New("path is empty/not found")
		}
	*/

	fileOrDir, err := chkFileOrDir(configFilePath)

	if err != nil {
		fmt.Printf("\nError: userConfig File/Path '%s' is not found", configFilePath)
		return err
	}

	configFileName := configFilePath

	if fileOrDir {
		configFileName = configFilePath + "/" + "config.sorl"
	}

	ok := true

	if !chkDir(configFileName) {
		fmt.Printf("\nUser specific Sorl config file '%s' not found.", configFileName)
		ok = false
	}

	//scProp := SorlConfigProperty{}

	if ok {
		fmt.Printf("\nDefault Sorl config file '%s' is found.", configFileName)
		scProp.readConfig(configFileName)
		//scProp.printConfig()
	}

	return nil
}
