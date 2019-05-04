package main

import (
	"fmt"
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

type SorlConfig map[string]string
type SorlConfigProperty map[string]SorlConfig

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

func (scProp SorlConfigProperty) printSection(matchStr, dispStr string) {

	for key, val := range scProp {
		if strings.HasPrefix(key, matchStr) {

			fmt.Printf("\n%s=%s", dispStr, key)
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
	fmt.Println("\n")

	scProp.printSection("sg:", "SuperGroups")
	fmt.Println("\n")

	scProp.printSection("g:", "GroupName")
	fmt.Println("\n")

	scProp.printSection("h:", "HostName")
	fmt.Println("\n" + strings.Repeat("=", cnt))
	fmt.Println("\n")
}

func (scProp SorlConfigProperty) readConfig(configFile string) {

	//lines, _ := ReadFile("./src/github.com/venugr/SOLcode/SOLutil/config.sorl")
	//PrintList(lines)
	lines, _ := ReadFile(configFile)

	hostName := ""
	groupName := ""
	allGroups := ""
	var sConfig, grpConfig SorlConfig

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
			}

			sConfig[tagKey] = tagVal
			scProp["h:"+hostName] = sConfig

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

		if strings.HasPrefix(line, "sorl_group") {
			idx := strings.Index(line, "=")
			tagKey := line[:idx]
			tagVal := strings.TrimSpace(line[idx+1:])

			if strings.HasPrefix(line, "sorl_group_name=") {
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

	}
}
