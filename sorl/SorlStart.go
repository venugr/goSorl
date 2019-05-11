package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var wgOrch = sync.WaitGroup{}

func sorlStart(orchFile string, scProp SorlConfigProperty, hostsList []string) {

	varsPerHostMap := make([]SorlMap, len(hostsList))

	for idx, lHost := range hostsList {

		lHostConfig := scProp["h:"+lHost]
		fmt.Printf("\nHost: %s", lHost)
		for lKey, lVal := range lHostConfig {
			fmt.Printf("\n\t%s=%s", lKey, lVal)

			if strings.EqualFold(lKey, "sorl_host_vars_file") {
				lVarsMap := SorlMap{}
				readVarsFile(lVal, &lVarsMap)
				printMap("Vars File: "+lVal, lVarsMap)
				varsPerHostMap[idx] = lVarsMap
			}
		}

	}

	for _, lHost := range hostsList {
		wgOrch.Add(1)
		go sorlProcessOrchestration(orchFile, lHost, scProp)
		wgOrch.Wait()

	}

	//wgOrch.Wait()

}

func sorlProcessOrchestration(orchFile, lHost string, scProp SorlConfigProperty) error {

	varsPerHostMap := SorlMap{}
	lHostConfig := scProp["h:"+lHost]
	lHostUser := ""
	lHostUserPasswd := ""
	lHostIP := ""
	lHostPort := 22

	fmt.Println("\n\ninfo: processing orchestration for Host:", lHost)

	for lKey, lVal := range lHostConfig {
		fmt.Printf("\t%s=%s\n", lKey, lVal)

		if strings.EqualFold(lKey, "sorl_host_vars_file") {
			readVarsFile(lVal, &varsPerHostMap)
			printMap("Vars File: "+lVal, varsPerHostMap)

		}

		if strings.EqualFold(lKey, "sorl_host_orch_file") {
			orchFile = lVal
		}

	}

	if lVal, ok := lHostConfig["sorl_host_vars_file"]; ok {
		readVarsFile(lVal, &varsPerHostMap)
		printMap("Vars File: "+lVal, varsPerHostMap)
	}

	if lVal, ok := lHostConfig["sorl_host_orch_file"]; ok {
		orchFile = lVal
	}

	if lVal, ok := lHostConfig["sorl_host_ssh_port"]; ok {
		lPort, err := strconv.Atoi(lVal)
		if err != nil {
			fmt.Printf("error: in valid port num: %v for the host: %s ", lVal, lHost)
			return err
		}

		lHostPort = lPort
	}

	if lVal, ok := lHostConfig["sorl_host_user"]; ok {
		lHostUser = lVal
	}

	if lVal, ok := lHostConfig["sorl_host_user_pass"]; ok {
		lHostUserPasswd = lVal
	}

	if lVal, ok := lHostConfig["sorl_host_ip"]; ok {
		lHostIP = lVal
	}

	session, client, err := sorlParallerlSsh(lHostUser, lHostUserPasswd, lHostIP, lHostPort)

	if err != nil {
		fmt.Printf("\nerror: session is not created due to: %v", err)
		fmt.Printf("\nerror: unable to proceed with orchestration")
	}

	runShell(session)
	defer session.Close()
	defer client.Close()

	//sorlStartOrchestration(orchFile, lHost, varsPerHostMap, scProp)
	//time.Sleep(2 * time.Second)
	wgOrch.Done()

	return nil
}
