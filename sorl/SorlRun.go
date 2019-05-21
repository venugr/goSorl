package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func sorlRunOrchestration(session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {

	//fmt.Println("Run Orchestration....")

	orchFile := string((*allProp)["sr:orchfile"])
	//color := string(allProp["sr:color"])
	keepNoCmdLogs, _ := strconv.Atoi((*allProp)["sr:keep"])
	keepCmdLogs := make([]string, keepNoCmdLogs)
	loadFile := string((*allProp)["sr:loadfile"])
	loadOk := string((*allProp)["sr:load"])
	//display := string(allProp["sr:display"])

	if loadOk == "yes" {
		orchFile = loadFile
	}

	commands, _ := ReadFile(orchFile)

	cmdOut := ""
	prevWaitCmd := ""
	runWaitOk := false

	if loadOk == "no" {
		//waitFor(color, []string{"$", "[BAN83] ?"}, sshIn)
	}
	for _, cmd := range commands {

		if checkPauseAbort() {
			fmt.Println("info: abort file is found")
			fmt.Println("info: aborting orchestration")
			return
		}

		cmd, err1 := replaceProp(cmd, Property(*allProp))

		cmd = strings.TrimLeft(cmd, " ")

		if runWaitOk && (!strings.HasPrefix(cmd, "wait")) {
			_, cmdOut = sorlOrchWait(prevWaitCmd, session, sshIn, sshOut, allProp)
		}
		runWaitOk = false

		//prevWaitCmd = cmd

		if strings.HasPrefix(cmd, "var ") {
			sorlOrchVar(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		if strings.HasPrefix(cmd, "load ") {
			sorlOrchLoad(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		cmdOut = ""
		if strings.HasPrefix(cmd, "wait ") {
			_, cmdOut = sorlOrchWait(cmd, session, sshIn, sshOut, allProp)
			prevWaitCmd = cmd
			runWaitOk = false
			continue
		}

		checkError(err1)
		runWaitOk = true
		//color := (*allProp)["sr:color"]
		//display := (*allProp)["sr:display"]
		//sshPrint(color, cmd+"\n")
		runShellCmd(cmd, sshOut)
		//if cmd != "exit" {
		//_, cmdOut := waitFor(color, []string{"$"}, sshIn)
		//}

		for i := 0; i < keepNoCmdLogs-1; i++ {
			keepCmdLogs[i] = keepCmdLogs[i+1]
		}
		keepCmdLogs[keepNoCmdLogs-1] = cmdOut

	}

	//session.Wait()
}

func sorlOrchVar(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {
	cmd = strings.Replace(cmd, "var ", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	vars := strings.Split(cmd, "=")
	(*allProp)[vars[0]] = vars[1]

	//printMap("Var Map", SorlMap(allProp))
}

func sorlOrchWait(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (int, string) {

	cmd = strings.Replace(cmd, "wait", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	waitStr := strings.Split(cmd, "||")
	color := (*allProp)["sr:color"]
	display := (*allProp)["sr:display"]

	return waitFor(color, display, waitStr, sshIn)

}

func sorlOrchLoad(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {

	loadFile := strings.Split(cmd, " ")

	locProp := Property{}

	for lKey, lVal := range *allProp {
		locProp[lKey] = lVal
	}

	locProp["sr:orchfile"] = loadFile[1]
	locProp["sr:load"] = "yes"
	//fmt.Println("Loading...", loadFile[1])

	(*allProp)["sr:loadfile"] = loadFile[1]
	(*allProp)["sr:load"] = "yes"
	sorlRunOrchestration(session, sshIn, sshOut, allProp)
}

func checkPauseAbort() bool {

	pauseFile := "/tmp/.sorl/.pause.sorl"
	abortFile := "/tmp/.sorl/.abort.sorl"
	ok := true

	if chkFile(abortFile) {
		return true
	}

	for {
		if chkFile(pauseFile) {
			if ok {
				fmt.Println("\nPause file is found")
				fmt.Println("Orchestration is paused")
				ok = false
			} else {
				fmt.Println("\nOrchestration is still paused")
			}
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if !ok {
		fmt.Println("\nOrchestration is resumed")
	}

	return false
}
