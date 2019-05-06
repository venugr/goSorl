package main

import (
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

	fmt.Println("\n")
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

	return cliArgs

}
