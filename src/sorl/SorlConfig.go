package main

import (
	"fmt"
	"os"
	"strings"
)

/*
type SorlConfigOld struct {
	solr_host_name          string
	solr_host_ip            string
	sorl_host_user          string
	sorl_host_user_pass     string
	sorl_host_ssh_keys      string
	sorl_host_ssh_keys_file string
	sorl_host_ssh_pass_less string
	sorl_host_group         string
}
*/

/*
func main() {

	scProp := SorlConfigProperty{}

	(&scProp).readConfig()

	for key, val := range scProp {
		if strings.HasPrefix(key, "g:") {
			fmt.Printf("\n\nGroupName=%s", key)
		} else if strings.HasPrefix(key, "ag:") {
			fmt.Printf("\n\nAllGroups=%s", key)
		} else {
			fmt.Printf("\n\nHostName=%s", key)
		}

		for k, v := range val {
			fmt.Printf("\n\t%s=%s", k, v)

		}
	}

	fmt.Println("\n")
}
*/

func sorlLoadConfigFiles(scProp *SorlConfigProperty, homePath string, userConfigFilePath string) {

	sorlDefaultConfigFile := homePath + PathSep + ".sorl" + PathSep + "config.sorl"
	readConfigFilePath(scProp, sorlDefaultConfigFile)

	//userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	if userConfigFilePath != "" {
		if err := readConfigFilePath(scProp, userConfigFilePath); err != nil {
			fmt.Printf("Unable to read proceed, %v", err)
			os.Exit(1)
		}
	}

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
	if sorlDebug {
		fmt.Printf("\nDefault Sorl config file '%s' is found.", configFileName)
	}

	if ok {
		scProp.readConfig(configFileName)
		//scProp.printConfig()
	}

	return nil
}

func (scProp SorlConfigProperty) printSection(matchStr, dispStr string) {

	for key, val := range scProp {
		if strings.HasPrefix(key, matchStr) {

			fmt.Printf("\n\n%s=%s", dispStr, key)
			for k, v := range val {
				fmt.Printf("\n\t%s=%s", k, v)

			}
		}
	}
}

func (scProp SorlConfigProperty) printConfig() {

	fmt.Println("\n" + strings.Repeat("=", cnt))
	fmt.Println("\t\tConfiguration Details")
	fmt.Println("\n" + strings.Repeat("=", cnt))
	scProp.printSection("ag:", "AllGroups")
	//fmt.Println("\n")

	scProp.printSection("sg:", "SuperGroups")
	//fmt.Println("\n")

	scProp.printSection("g:", "GroupName")
	//fmt.Println("\n")

	scProp.printSection("h:", "HostName")

	fmt.Println("\n" + strings.Repeat("=", cnt))
	fmt.Println()
	fmt.Println()
}

func (scProp SorlConfigProperty) readConfig(configFile string) {

	//lines, _ := ReadFile("./src/github.com/venugr/SOLcode/SOLutil/config.sorl")
	//PrintList(lines)
	lines, _ := ReadFile(configFile)

	hostName := ""
	groupName := ""
	allGroups := ""
	var sConfig, logConfig, grpConfig, allHostNames, allGroup2Hosts SorlConfig

	allHostNames = SorlConfig{}
	if val, ok := scProp["all.hosts"]; ok {
		allHostNames = val
	}
	scProp["all.hosts"] = allHostNames

	allGroup2Hosts = SorlConfig{}
	if val, ok := scProp["group.hosts"]; ok {
		allGroup2Hosts = val
	}
	scProp["group.hosts"] = allGroup2Hosts

	for _, line := range lines {

		//fmt.Println(line)
		if strings.HasPrefix(line, "sorl_host") {
			idx := strings.Index(line, "=")
			tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])
			//fmt.Println(hostName)
			//sConfig := SorlConfig{}
			//sConfig[tagKey] = tagVal

			if strings.HasPrefix(line, "sorl_host_name=") {
				hostName = tagVal
				sConfig = SorlConfig{}
				allHostNames[hostName] = "NA"
				scProp["all.hosts"] = allHostNames

			}

			sConfig[tagKey] = tagVal
			scProp["h:"+hostName] = sConfig

			if strings.HasPrefix(line, "sorl_host_group=") {

				mapKV := scProp["group.hosts"]
				for _, grp := range strings.Split(tagVal, ",") {
					if mVal, ok := mapKV[grp]; ok {
						mVal = mVal + "," + hostName
						mapKV[grp] = mVal
						scProp["group.hosts"] = mapKV
					} else {
						mapKV[grp] = hostName
						scProp["group.hosts"] = mapKV
					}

				}

			}

			//fmt.Println(sConfig)

		}

		if strings.HasPrefix(line, "sorl_host_group=") {
			idx := strings.Index(line, "=")
			//tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])
			hostGroups := strings.Split(tagVal, ",")

			for _, grp := range hostGroups {

				if val, ok := scProp["g:"+grp]; ok {
					val["hosts"] = val["hosts"] + "," + hostName
				} else {
					grpConfig = SorlConfig{}
					grpConfig["hosts"] = hostName
					scProp["g:"+grp] = grpConfig
				}

			}

		}

		if strings.HasPrefix(line, "sorl_group_") {
			idx := strings.Index(line, "=")
			tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])

			if strings.HasPrefix(line, "sorl_group_super_name=") {
				groupName = tagVal
				sConfig = SorlConfig{}
			}

			sConfig[tagKey] = tagVal
			scProp["sg:"+groupName] = sConfig

			//fmt.Println(sConfig)

		}

		if strings.HasPrefix(line, "sorl_all_groups") {
			idx := strings.Index(line, "=")
			tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])

			if strings.HasPrefix(line, "sorl_all_groups=") {
				allGroups = tagVal
				sConfig = SorlConfig{}
			}

			sConfig[tagKey] = tagVal
			scProp["ag:"+allGroups] = sConfig

			//fmt.Println(sConfig)

		}

		if strings.HasPrefix(line, "sorl_log_path") {
			idx := strings.Index(line, "=")
			tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])
			logConfig = SorlConfig{}
			logConfig[tagKey] = tagVal
			scProp["lp:logpath"] = logConfig

			//fmt.Println(sConfig)

		}

	}
}
