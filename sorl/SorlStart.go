package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wgOrch = sync.WaitGroup{}

func sorlStart(parallelOk, orchFile string, scProp SorlConfigProperty,
	hostsList []string, cliArgsMap map[string]string,
	svMap map[string]string) {

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
		"\x1b[36;1m",
	}
	max := len(SorlColors)
	min := 1
	rand.Seed(time.Now().UnixNano())

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
		trand := rand.Intn(max-min) + min
		go sorlProcessOrchestration(SorlColors[trand], orchFile, lHost, scProp, cliArgsMap, svMap)

		if parallelOk == "false" {
			wgOrch.Wait()
		}

	}

	if parallelOk == "true" {
		wgOrch.Wait()
	}

}

func sorlProcessOrchestration(color, orchFile, lHost string, scProp SorlConfigProperty,
	cliArgsMap map[string]string, svMap map[string]string) error {

	varsPerHostMap := SorlMap{}
	lHostConfig := scProp["h:"+lHost]
	lLogConfig := scProp["lp:logpath"]
	lHostUser := ""
	lHostUserPasswd := ""
	lHostIP := ""
	lHostPort := 22
	lHostLogName := "_sorl.log"
	keepNoCmdLogs, _ := strconv.Atoi(cliArgsMap["keep"])
	keepCmdLogs := make([]string, keepNoCmdLogs)
	lLogPath := lLogConfig["sorl_log_path"]

	for fKey, fVal := range svMap {
		varsPerHostMap[fKey] = fVal
	}

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

	if lVal, ok := lHostConfig["sorl_host_name"]; ok {
		lHostLogName = lVal + lHostLogName
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

	lLogPath = lLogPath + PathSep + lHostLogName

	fmt.Println("Orchestration file:", orchFile)
	fmt.Println("          Log file:", lLogPath)

	time.Sleep(2 * time.Second)

	session, client, err := sorlParallelSsh(lHostUser, lHostUserPasswd, lHostIP, lHostPort)

	if err != nil {
		fmt.Printf("\nerror: session is not created due to: %v", err)
		fmt.Printf("\nerror: unable to proceed with orchestration")
	}

	//runShell(session)

	sshIn, sshOut, sshErr := setShell(session)

	if sshErr != nil {
		fmt.Println(sshErr)
		os.Exit(1)
	}

	/*
		commands := []string{
			"uname -a",
			"pwd",
			"sleep 5",
			"pwd",
			"echo 'bye'",
			"ls -l /tmp",
			"sqlplus --h",
			"sleep 2",
			"env | sort ",
			"ls",
			"#sqlplus /nolog @/media/common/db/versions",
			"df -h",
			"exit",
		}
	*/

	commands, _ := ReadFile(orchFile)

	//PrintList("FILE", commands)
	waitFor(color, []string{"$", "[BAN83] ?"}, sshIn)
	for _, cmd := range commands {
		cmd, err1 := replaceProp(cmd, Property(varsPerHostMap))
		checkError(err1)
		runShellCmd(cmd, sshOut)
		//if cmd != "exit" {
		_, cmdOut := waitFor(color, []string{"$"}, sshIn)
		//}

		for i := 0; i < keepNoCmdLogs-1; i++ {
			keepCmdLogs[i] = keepCmdLogs[i+1]
		}
		keepCmdLogs[keepNoCmdLogs-1] = cmdOut

	}

	session.Wait()
	defer session.Close()
	defer client.Close()

	if false {
		for i := 0; i < keepNoCmdLogs; i++ {
			fmt.Printf("\n\nLog[%v]:%s", i+1, keepCmdLogs[i])
		}
	}
	//sorlStartOrchestration(orchFile, lHost, varsPerHostMap, scProp)
	//time.Sleep(2 * time.Second)
	wgOrch.Done()

	return nil
}
