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
	infoStr := ""

	for _, key := range keys {
		val, ok := getEnvVal(key)
		infoStr = "Environment variable " + key + " is not available"
		if ok {
			infoStr = "Environment variable " + key + ":" + val
			keyVal[key] = val
		}

		fmt.Println(infoStr)

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

	connectToPtr := flag.String("connect-to", "", "Connect to system")
	connectUserPtr := flag.String("connect-user", "", "Connect as User")
	connectPasswordEncPtr := flag.String("connect-password-enc", "", "Encrypted Password")
	connectAskPsswordPtr := flag.Bool("connect-ask-password", false, "Ask for Password")
	connectPasswordKeysFilePtr := flag.String("connect-passwordless-keys-file", "", "Passwordless keys file path")
	flag.Parse()

	/*
		fmt.Println("Group:", *groupPtr)
		fmt.Println("Max Go Routines:", *maxGoPtr)
		fmt.Println("Dry Run:", *dryRunPtr)
		fmt.Println("Host Name:", *hostNamePtr)
		fmt.Println("Tail:", flag.Args())
	*/

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
	cliArgs["connect-user"] = string(*connectUserPtr)
	cliArgs["connect-password-enc"] = string(*connectPasswordEncPtr)
	cliArgs["connect-passwordless-keys-file"] = string(*connectPasswordKeysFilePtr)

	cliArgs["connect-ask-password"] = "false"
	if *connectAskPsswordPtr {
		cliArgs["connect-ask-password"] = "true"
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

func sorlConnectCliArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) bool {

	connectToCli := strings.TrimSpace(cliArgsMap["connect-to"])

	fmt.Println("ConnectTo:" + connectToCli)

	if connectToCli == "" {
		return false
	}

	connectUser := strings.TrimSpace(cliArgsMap["connect-user"])

	if connectUser == "" {
		connectUsage()
		return true
	}

	connectPasswordEnc := strings.TrimSpace(cliArgsMap["connect-password-enc"])
	connectPasswordlessKeysFile := strings.TrimSpace(cliArgsMap["connect-passwordless-keys-file"])
	connectAskPassword := strings.TrimSpace(cliArgsMap["connect-ask-password"])

	if connectPasswordEnc == "" &&
		connectPasswordlessKeysFile == "" &&
		connectAskPassword == "false" {

		connectUsage()
		return true
	}

	return true
}

func connectUsage() {

	fmt.Println("\nError: One of the following arguments needed for '-connect-to' ")
	fmt.Println("\t --connect-user")
	fmt.Println("\t --connect-password-enc")
	fmt.Println("\t --connect-passwordless-keys-file")
	fmt.Println("\t --connect-ask-password")

	fmt.Println()

}

func sorlProcessCliArgs(scProp SorlConfigProperty, cliArgsMap map[string]string) ([]string, error) {

	//grpCliOk := false
	//hostCliOk := false

	hostCli := strings.TrimSpace(cliArgsMap["host"])
	grpCli := strings.TrimSpace(cliArgsMap["group"])
	//maxGoRout := cliArgsMap["max"]

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
