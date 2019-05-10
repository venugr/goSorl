package main

import (
	"fmt"
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

func sorlProcessOrchestration(orchFile, lHost string, scProp SorlConfigProperty) {

	varsPerHostMap := SorlMap{}
	lHostConfig := scProp["h:"+lHost]

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

	sorlStartOrchestration(orchFile, lHost, varsPerHostMap, scProp)
	//time.Sleep(2 * time.Second)
	wgOrch.Done()
}
