package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

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

	SorlColors := []string{
		"",
		"\x1b[31;1m",
		"\x1b[37;1m",
		"\x1b[32;1m",
		"\x1b[33;1m",
		"\x1b[34;1m",
		"\x1b[35;1m",
		"\x1b[36;1m",
		"\x1b[38;1m",
		"\x1b[39;1m",
	}
	max := len(SorlColors)
	min := 0
	rand.Seed(time.Now().UnixNano())
	trand := rand.Intn(max-min) + min
	color := SorlColors[trand]

	ss := &SorlSSH{}
	allProp := Property{}
	allProp["sr:debug"] = cliArgsMap["debug"]
	allProp["sr:info"] = cliArgsMap["info"]
	allProp["sr:echo"] = "on"
	allProp["sr:sorl_user_homepath"] = cliArgsMap["sorl_user_homepath"]
	allProp["sr:orchfile"] = scriptName
	//ss.sorlOrchestration(strings.Join(fileInfo, "\n"), &alp)
	//ss.sorlRunOrchestration(&allProp)

	if cliArgsMap["connect-to"] != "" {
		ss.sorlSshHostIP = cliArgsMap["connect-to"]
	}

	if cliArgsMap["connect-to"] == "" {
		ss.sorlRunOrchestration(&allProp)
		return nil
	}

	if cliArgsMap["conn-user"] != "" {
		ss.sorlSshUserName = cliArgsMap["conn-user"]
	}

	if cliArgsMap["conn-password-enc"] != "" {
		lTempPasswd := cliArgsMap["conn-password-enc"]
		if strings.HasPrefix(lTempPasswd, "sorl.enc:") {
			key := "123456789012345678901234"
			encPasswd := strings.Split(lTempPasswd, ".enc:")[1]
			lTempPasswd = sorlDecryptText(key, encPasswd)
		}
		ss.sorlSshUserPassword = lTempPasswd
	}

	if cliArgsMap["conn-port"] != "" {
		ss.sorlSshHostPortNum, _ = strconv.Atoi(cliArgsMap["conn-port"])
	}

	if cliArgsMap["conn-passwordless-keys-file"] != "" {
		ss.sorlSshHostKeyFile = cliArgsMap["conn-passwordless-keys-file"]
	}
	//ss.sorlSshColor = color

	err := ss.sorlParallelSsh()
	if err != nil {
		fmt.Printf("\nerror: session is not created due to: %v", err)
		fmt.Printf("\nerror: unable to proceed with orchestration")
	}

	sshErr := ss.setShell()
	if sshErr != nil {
		fmt.Println(sshErr)
		os.Exit(1)
	}

	if strings.HasPrefix(scriptName, "~") {
		scriptName = strings.Replace(scriptName, "~", cliArgsMap["sorl_user_homepath"], 1)
	}

	keepNoCmdLogs, _ := strconv.Atoi(cliArgsMap["keep"])
	keepCmdLogs := make([]string, keepNoCmdLogs)
	display, _ := cliArgsMap["display"]
	tags, _ := cliArgsMap["tags"]

	allProp["sr:load"] = "no"
	allProp["sr:loadfile"] = ""
	allProp["sr:orchfile"] = scriptName
	allProp["sr:color"] = color
	allProp["sr:keep"] = strconv.Itoa(keepNoCmdLogs)
	allProp["sr:display"] = display
	allProp["sr:tags"] = tags
	allProp["sr:debug"] = cliArgsMap["debug"]
	allProp["sr:info"] = cliArgsMap["info"]
	allProp["_host.local"] = "no"
	allProp["sr:echo"] = "on"

	ss.sorlRunOrchestration(&allProp)
	ss.sorlSshSession.Wait()

	defer ss.sorlSshSession.Close()
	defer ss.sorlSshClient.Close()

	if false {
		for i := 0; i < keepNoCmdLogs; i++ {
			fmt.Printf("\n\nLog[%v]:%s", i+1, keepCmdLogs[i])
		}
	}

	fmt.Println("\n\ninfo: processed orchestration for Host:", ss.sorlSshHostIP)
	//wgOrch.Done()

	return nil

}
