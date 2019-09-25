package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func sorlActionConn(actName, connSystem string, actArgs []string, cliArgsMap map[string]string, svMap SorlMap, scProp SorlConfigProperty) {

	if sorlDebug {
		fmt.Println("Action: " + actName)
		fmt.Println("Conn System: " + connSystem)
		fmt.Println("Conn Port: " + actArgs[0])
		fmt.Println("Conn User: " + actArgs[1])
		fmt.Println("Conn Pass: " + actArgs[2])
		fmt.Println("Conn PassFile: " + actArgs[3])
		fmt.Println("Conn Ask Passwd: " + actArgs[4])
		fmt.Println("Conn Cmds File: " + actArgs[6])

	}

	color := SorlGetColor()

	allProp := Property{}
	allProp["sr:load"] = "no"
	allProp["sr:loadfile"] = ""
	allProp["sr:orchfile"] = ""
	allProp["sr:color"] = color
	allProp["_host.local"] = "no"
	//allProp["sr:keep"] = strconv.Itoa(keepNoCmdLogs)
	//allProp["sr:display"] = display
	//allProp["sr:tags"] = tags
	//allProp["sr:debug"] = cliArgsMap["debug"]

	for fKey, fVal := range svMap {
		allProp[fKey] = fVal
	}

	ss := &SorlSSH{}
	ss.sorlSshHostIP = connSystem
	ss.sorlSshUserName = actArgs[1]
	ss.sorlSshUserPassword = actArgs[2]
	ss.sorlSshHostPortNum, _ = strconv.Atoi(actArgs[0])
	ss.sorlSshHostKeyFile = actArgs[3]
	ss.sorlSshColor = color
	waitPrompt := actArgs[5]
	connectCmdsFile := actArgs[6]

	if actArgs[4] == "true" {
		//reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n\nEnter password: ")
		//text, _ := reader.ReadString('\n')
		//text = strings.TrimRight(text, "\n")
		//ss.sorlSshUserPassword = text

		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err == nil {
			//fmt.Println("\nPassword typed: " + string(bytePassword))
		}
		password := string(bytePassword)
		ss.sorlSshUserPassword = strings.TrimSpace(password)

		fmt.Println()
	}

	lHostUserPasswd := ss.sorlSshUserPassword

	if strings.HasPrefix(lHostUserPasswd, "sorl.enc:") || strings.HasPrefix(lHostUserPasswd, "enc:") {
		key := "123456789012345678901234"
		encPasswd := strings.Split(lHostUserPasswd, "enc:")[1]
		lHostUserPasswd = sorlDecryptText(key, encPasswd)
	}
	ss.sorlSshUserPassword = lHostUserPasswd

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

	cmdStr := ""

	for _, lVal := range strings.Split(allProp["_cmd.arg.order"], ",") {

		cmdStr += allProp[lVal] + "\n"

	}

	if waitPrompt == "" {
		waitPrompt = "$||#||?||>"
	}

	waitPrompt = ".wait " + waitPrompt

	cmdStr = strings.TrimSpace(cmdStr)

	//fmt.Println("WaitPrompt: " + waitPrompt)

	commands := []string{
		waitPrompt,
		cmdStr,
		"exit",
	}

	if connectCmdsFile != "" {
		fileData := strings.Join(commands, "\n") + "\n"
		if ok := WriteFile(connectCmdsFile, fileData); ok != nil {
			fmt.Print("\nWarn: ")
			fmt.Println(ok)
		}
	}

	fmt.Println()

	ss.sorlOrchestration(strings.Join(commands, "\n"), &allProp)
	//ss.sorlSshSession.Wait()

	//defer ss.sorlSshSession.Close()
	//defer ss.sorlSshClient.Close()

	fmt.Println("\n\ninfo: processed orchestration for Host:", connSystem)
	//wgOrch.Done()

}
