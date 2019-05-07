package main

import "fmt"

func sorlStart(scProp SorlConfigProperty, hostsList []string) {

	for _, lHost := range hostsList {

		lHostConfig := scProp["h:"+lHost]
		fmt.Printf("\nHost: %s", lHost)
		for lKey, lVal := range lHostConfig {
			fmt.Printf("\n\t%s=%s", lKey, lVal)
		}

	}
}
