package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getEnvVal(key string) (string, bool) {

	return os.LookupEnv(key)
}

func getEnvlist(keys []string) map[string]string {

	keyVal := map[string]string{}
	//infoStr := ""

	for _, key := range keys {
		val, ok := getEnvVal(key)
		//infoStr = "Environment variable " + key + " is not available"
		if ok {
			//infoStr = "Environment variable " + key + ":" + val
			keyVal[key] = val
		}

		//fmt.Println(infoStr)

	}
	return keyVal
}

func getEnvVars() {

	fmt.Println()
	fmt.Println()

	for _, ev := range os.Environ() {
		pairs := strings.Split(ev, "=")
		fmt.Printf("%s=%s\n", pairs[0], pairs[1])
	}
}

func getCliArgs() map[string]string {

	cliArgs := map[string]string{}

	connPtr := flag.String("conn", "", "Connect to System")
	groupPtr := flag.String("group", "", "Group Name")
	maxGoPtr := flag.Int("max", 10, "Maximum Parallel Go Routines")
	dryRunPtr := flag.Bool("dryrun", false, "Dry Run")
	hostNamePtr := flag.String("host", "", "Host Name")
	configFilePtr := flag.String("config", "", "Config File Path")
	orchFilePtr := flag.String("orchfile", "", "Orchestration File Path")
	allowPtr := flag.Bool("parallel", false, "Allow Parallel Go Routines")
	keepPtr := flag.Int("keep", 5, "Keep no of command logs")
	dispPtr := flag.String("display", "more", "display: [less | more | all | no | clear]")
	tagsPtr := flag.String("tags", "", "Tag Names")
	versionPtr := flag.Bool("version", false, "Version details")
	debugPtr := flag.Bool("debug", false, "Include debug blocks")
	encyPtr := flag.String("encrypt", "", "Encrypt the string")
	decyPtr := flag.String("decrypt", "", "Decrypt the string")
	varPtr := flag.String("var", "", "Variable Name")
	varFilePtr := flag.String("var-file", "", "Variables FileName")
	cmdPtr := flag.String("cmd", "", "Command text")

	connectToPtr := flag.String("connect-to", "", "Connect to system")
	connectUserPtr := flag.String("conn-user", "", "Connect as User")
	connectPasswordEncPtr := flag.String("conn-password-enc", "", "Encrypted Password")
	connectAskPsswordPtr := flag.Bool("conn-ask-password", false, "Ask for Password")
	connectPasswordKeysFilePtr := flag.String("conn-passwordless-keys-file", "", "Passwordless keys file path")
	flag.Parse()

	/*
		fmt.Println("Group:", *groupPtr)
		fmt.Println("Max Go Routines:", *maxGoPtr)
		fmt.Println("Dry Run:", *dryRunPtr)
		fmt.Println("Host Name:", *hostNamePtr)
		fmt.Println("Tail:", flag.Args())
	*/

	cliArgs["conn"] = string(*connPtr)
	cliArgs["group"] = string(*groupPtr)
	cliArgs["max"] = strconv.Itoa(*maxGoPtr)
	cliArgs["dryrun"] = "false"
	if *dryRunPtr {
		cliArgs["dryrun"] = "true"
	}
	cliArgs["host"] = string(*hostNamePtr)
	cliArgs["config"] = string(*configFilePtr)
	cliArgs["orchfile"] = string(*orchFilePtr)
	cliArgs["parallel"] = "false"
	if *allowPtr {
		cliArgs["parallel"] = "true"
	}

	cliArgs["keep"] = strconv.Itoa(*keepPtr)
	cliArgs["display"] = string(*dispPtr)
	cliArgs["tags"] = string(*tagsPtr)
	cliArgs["var"] = string(*varPtr)
	cliArgs["var-file"] = string(*varFilePtr)
	cliArgs["cmd"] = string(*cmdPtr)

	cliArgs["version"] = "false"
	if *versionPtr {
		cliArgs["version"] = "true"
	}

	cliArgs["debug"] = "false"
	if *debugPtr {
		cliArgs["debug"] = "true"
	}

	cliArgs["encrypt"] = string(*encyPtr)
	cliArgs["decrypt"] = string(*decyPtr)

	cliArgs["connect-to"] = string(*connectToPtr)
	cliArgs["conn-user"] = string(*connectUserPtr)
	cliArgs["conn-password-enc"] = string(*connectPasswordEncPtr)
	cliArgs["conn-passwordless-keys-file"] = string(*connectPasswordKeysFilePtr)

	cliArgs["conn-ask-password"] = "false"
	if *connectAskPsswordPtr {
		cliArgs["conn-ask-password"] = "true"
	}

	return cliArgs

}

func sorlEncDecCliArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) bool {

	encyCli := strings.TrimSpace(cliArgsMap["encrypt"])
	decyCli := strings.TrimSpace(cliArgsMap["decrypt"])

	if encyCli != "" && decyCli != "" {
		fmt.Println("\nError: Both 'encrypt' and 'decrypt' can not be present.")
		fmt.Println()
		os.Exit(1)
	}

	if encyCli != "" {
		key := "123456789012345678901234"
		encStr := sorlEncryptText(key, encyCli)
		fmt.Println("String:" + encyCli)
		fmt.Println("sorl.enc:" + encStr)
		return true
	}

	if decyCli != "" {
		key := "123456789012345678901234"
		decStr := sorlDecryptText(key, decyCli)
		fmt.Println("String:" + decyCli)
		fmt.Println("sorl.dec:" + decStr)
		return true
	}

	return false
}

func sorlConnectCliArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) ([]string, error) {

	/*
		connectToCli := strings.TrimSpace(cliArgsMap["connect-to"])

		fmt.Println("ConnectTo:" + connectToCli)

		if connectToCli == "" {
			return false
		}
	*/

	connectUser := strings.TrimSpace(cliArgsMap["conn-user"])

	if connectUser == "" {
		connectUsage()
		return nil, errors.New("Insufficient Arguments")
	}

	connectPasswordEnc := strings.TrimSpace(cliArgsMap["conn-password-enc"])
	connectPasswordlessKeysFile := strings.TrimSpace(cliArgsMap["conn-passwordless-keys-file"])
	connectAskPassword := strings.TrimSpace(cliArgsMap["conn-ask-password"])

	if connectPasswordEnc == "" &&
		connectPasswordlessKeysFile == "" &&
		connectAskPassword == "false" {

		connectUsage()
		return nil, errors.New("Insufficient Arguments")
	}

	connCnt := 0

	if connectPasswordEnc != "" {
		connCnt++
	}

	if connectPasswordlessKeysFile != "" {
		connCnt++
	}

	if connectAskPassword != "false" {
		connCnt++
	}

	if connCnt == 0 || connCnt > 1 {
		//connectUsage()
		return nil, errors.New("Error: More than one password argumnet is present\n")
	}

	//sorlStart(parallelOk, globalOrchFilePath, scProp, hostList, cliArgsMap, svMap)

	return []string{connectUser, connectPasswordEnc, connectPasswordlessKeysFile, connectAskPassword}, nil
}

func connectUsage() {

	fmt.Println("\nError: Following arguments are required for '-conn' ")
	fmt.Println("\t   --conn-user")
	fmt.Println("\t        AND")
	fmt.Println("\t [ --conn-password-enc")
	fmt.Println("\t        OR")
	fmt.Println("\t   --conn-passwordless-keys-file")
	fmt.Println("\t        OR")
	fmt.Println("\t   --conn-ask-password")
	fmt.Println("\t ]")

	fmt.Println()

}

func sorlGetActionArgs(actName string, scProp SorlConfigProperty, cliArgsMap map[string]string) ([]string, error) {

	if actName == "host" || actName == "group" {
		return getHostList(actName, strings.TrimSpace(cliArgsMap[actName]), scProp)
	}

	if actName == "conn" {
		return sorlConnectCliArgs(scProp, cliArgsMap)
	}

	return nil, nil
}

func sorlGetAction(cliArgsMap map[string]string) (string, error) {

	actList := []string{"host", "group", "conn", "encrypt", "decrypt"}
	actCnt := 0
	actName := ""

	for _, tVal := range actList {
		if strings.TrimSpace(cliArgsMap[tVal]) != "" {
			actCnt++
			actName = tVal
		}
	}

	if actCnt == 0 || actCnt > 1 {
		fmt.Println("\nError: Only one of the following Actions be present.")
		for _, tVal := range actList {
			fmt.Println("\t- " + tVal)
		}
		fmt.Println()

		return "", nil
	}

	return actName, nil

}

func sorlProcessCliArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) ([]string, error) {

	//grpCliOk := false
	//hostCliOk := false

	hostCli := strings.TrimSpace(cliArgsMap["host"])
	grpCli := strings.TrimSpace(cliArgsMap["group"])
	connCli := strings.TrimSpace(cliArgsMap["conn"])

	//maxGoRout := cliArgsMap["max"]

	if hostCli != "" && grpCli != "" && connCli != "" {
		fmt.Println("\nError: One of the actions 'conn', host' and 'group' can be present.")
		fmt.Println()
		os.Exit(1)
	}

	if hostCli == "" && grpCli == "" && connCli == "" {
		fmt.Println("\nError: One of 'conn', host' or 'group' must be present.")
		fmt.Println()
		os.Exit(1)
	}

	if connCli != "" {
		return []string{connCli}, nil
	}

	hostGrpCli := hostCli
	//hostCliOk = true
	selType := "host"

	if strings.TrimSpace(grpCli) != "" {
		hostGrpCli = grpCli
		selType = "group"
		//grpCliOk = true
		//hostCliOk = false
	}

	fmt.Println("Host/Group:", hostGrpCli)
	//fmt.Println(grpCliOk, hostCliOk, selType, maxGoRout)

	hostList, err := getHostList(selType, hostGrpCli, scProp)
	if err != nil {
		return nil, err
	}

	return hostList, nil

}

func getHostList(selType, hostGrpCli string, scProp SorlConfigProperty) ([]string, error) {
	hostList := []string{}

	if strings.EqualFold(selType, "host") && hostGrpCli != "all" {
		allHostNames := scProp["all.hosts"]
		allHostsStr := ""

		for k := range allHostNames {
			allHostsStr += k + ","
		}

		for _, lVal := range strings.Split(hostGrpCli, ",") {

			if !strings.Contains(allHostsStr, lVal+",") {
				fmt.Printf("\nerror: config info for host '%v' is not found", lVal)
				fmt.Printf("\nerror: aborting orchestration...")
				return nil, errors.New("host info not found")
			}
		}

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
				return nil, fmt.Errorf("Error: group id '%s' is not found", grpVal)
			}

		}
		hostsStr = strings.Trim(hostsStr, ",")

		return strings.Split(hostsStr, ","), nil

	}

	return hostList, nil
}
