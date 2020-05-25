package main

func sorlRunScript(scriptName string, scProp SorlConfigProperty, cliArgsMap map[string]string) error {

	/*
		if strings.HasPrefix(scriptName, "~/") {
			scriptName = strings.Replace(scriptName, "~", cliArgsMap["sorl_user_homepath"], 1)
		}

		fileInfo, err := ReadFile(scriptName)

		if err != nil {
			fmt.Println(err)
			return err
		}
	*/

	ss := &SorlSSH{}
	alp := Property{}
	alp["sr:debug"] = cliArgsMap["debug"]
	alp["sr:info"] = cliArgsMap["info"]
	alp["sr:echo"] = "on"
	alp["sr:sorl_user_homepath"] = cliArgsMap["sorl_user_homepath"]
	alp["sr:orchfile"] = scriptName
	//ss.sorlOrchestration(strings.Join(fileInfo, "\n"), &alp)
	ss.sorlRunOrchestration(&alp)

	return nil

}
