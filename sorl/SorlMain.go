package main

import "fmt"

func main() {

	fmt.Println()

	cliArgsMap := getCliArgs()
	fmt.Println(cliArgsMap)

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	fmt.Println(envMap)

	homePath := envMap["HOME"]
	sorlConfigFile := homePath + "/" + ".sorl" + "/" + "config.sorl"

	ok := true

	if !chkDir(sorlConfigFile) {
		fmt.Printf("\nSorl config file '%s' not found.", sorlConfigFile)
		ok = false
	}

	scProp := SorlConfigProperty{}

	if ok {
		(&scProp).readConfig(sorlConfigFile)
		scProp.printConfig()
	}

	fmt.Println()
}
