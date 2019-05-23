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
	(*allProp)["_cmd.output"] = ""

	orchFile := string((*allProp)["sr:orchfile"])
	//color := string((*allProp)["sr:color"])
	keepNoCmdLogs, _ := strconv.Atoi((*allProp)["sr:keep"])
	keepCmdLogs := make([]string, keepNoCmdLogs)
	loadFile := string((*allProp)["sr:loadfile"])
	loadOk := string((*allProp)["sr:load"])
	//display := string(allProp["sr:display"])

	skipTagLines := false
	skipIfLines := false
	//tagName := ""

	if loadOk == "yes" {
		orchFile = loadFile
	}

	commands, _ := ReadFile(orchFile)

	cmdOut := ""
	prevWaitCmd := ""
	runWaitOk := false
	tempCmdOut := ""
	skipVarLines := false
	isRemoved := false

	if loadOk == "no" {
		//waitFor(color, []string{"$", "[BAN83] ?"}, sshIn)
	}

	for _, cmd := range commands {

		if runWaitOk && (!strings.HasPrefix(cmd, ".wait")) {
			_, cmdOut = sorlOrchWait(prevWaitCmd, session, sshIn, sshOut, allProp)
			tempCmdOut += cmdOut
			(*allProp)["_cmd.output"] = tempCmdOut
		}
		runWaitOk = false

		if checkPauseAbort() {
			fmt.Println("info: abort file is found")
			fmt.Println("info: aborting orchestration")
			return
		}

		//cmd, err1 := replaceProp(cmd, Property(*allProp))
		cmd = strings.TrimLeft(cmd, " ")

		if strings.HasPrefix(cmd, "#") || cmd == "" {
			continue
		}

		if strings.HasPrefix(cmd, "}") {
			//fmt.Println("Skipping...1:", cmd)
			if skipTagLines {
				skipTagLines = false

			}
			if skipIfLines {
				skipIfLines = false
			}

			if skipVarLines {
				skipVarLines = false
			}

			continue
		}

		if skipVarLines {
			sorlOrchVar(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		if skipTagLines || skipIfLines {
			//fmt.Println("Skipping...:", cmd)
			continue
		}

		isRemoved = false
		if strings.HasSuffix(cmd, "{") {
			cmd = strings.TrimRight(cmd, "{")
			isRemoved = true
		}

		cmd, err1 := replaceProp(cmd, Property(*allProp))
		checkError(err1)

		if isRemoved {
			cmd = cmd + "{"
		}

		if strings.HasPrefix(cmd, ".sleep") {
			sorlOrchSleep(cmd)
			continue
		}

		if strings.HasPrefix(cmd, ".clear") {
			tempCmdOut = ""
			(*allProp)["_cmd.output"] = tempCmdOut
			continue
		}

		if strings.HasPrefix(cmd, ".print") {
			sorlOrchPrint(cmd, (*allProp)["sr:color"])
			continue
		}

		if strings.HasPrefix(cmd, ".pass") {
			if sorlOrchPass(cmd, (*allProp)["sr:color"], tempCmdOut) {
				sshPrint((*allProp)["sr:color"], "\n"+cmd+" : Failed\n")
				session.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(cmd, ".fail") {
			if sorlOrchFail(cmd, (*allProp)["sr:color"], tempCmdOut) {
				sshPrint((*allProp)["sr:color"], "\n"+cmd+" : Failed\n")
				session.Close()
				return
			}
			continue
		}

		//prevWaitCmd = cmd

		if strings.HasPrefix(cmd, ".tag ") {
			skipTagLines, _ = sorlOrchTag(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".if ") {
			skipIfLines, _ = sorlOrchIf(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".var ") && (!strings.Contains(cmd, "{")) {
			//fmt.Printf("var single %v", cmd)
			sorlOrchVar(cmd, session, sshIn, sshOut, allProp)
			//printMap("Var Map", SorlMap(*allProp))
			continue
		}

		if strings.HasPrefix(cmd, ".var ") && strings.Contains(cmd, "{") {
			skipVarLines = true
			//fmt.Printf("var group")
			continue
		}

		if strings.HasPrefix(cmd, ".load ") {
			sorlOrchLoad(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		cmdOut = ""
		if strings.HasPrefix(cmd, ".wait ") {
			_, cmdOut = sorlOrchWait(cmd, session, sshIn, sshOut, allProp)
			prevWaitCmd = cmd
			runWaitOk = false
			tempCmdOut += cmdOut
			(*allProp)["_cmd.output"] = tempCmdOut
			continue
		}

		if strings.HasPrefix(cmd, ".enter") {
			cmd = ""
		}

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

func sorlOrchFail(cmd, color, tempCmdOut string) bool {

	cmd = strings.Replace(cmd, ".fail ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		return true
	}

	return false

}

func sorlOrchPass(cmd, color, tempCmdOut string) bool {
	cmd = strings.Replace(cmd, ".pass ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		return false
	}

	return true

}

func getIfData(cmd, orStr, andStr, eqStr, nEqStr string) (string, string) {

	idxMap := map[string]int{}
	idx := -1
	condStr := ""

	idxMap[orStr] = strings.Index(cmd, orStr)
	idxMap[andStr] = strings.Index(cmd, andStr)
	idxMap[eqStr] = strings.Index(cmd, eqStr)
	idxMap[nEqStr] = strings.Index(cmd, nEqStr)

	for lKey, lVal := range idxMap {

		if lVal == -1 {
			continue
		}

		if idx == -1 && lVal > idx {
			idx = lVal
			condStr = lKey
			continue
		}

		if lVal < idx {
			idx = lVal
			condStr = lKey
		}

	}

	if idx == -1 {
		return cmd, ""
	}
	return cmd[:idx], condStr

}

func sorlOrchIf(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	prevVal1 := ""
	prevOp1 := ""
	tVal1 := ""
	tOp1 := ""
	//fmt.Println("inside...tag")
	cmd = strings.Replace(cmd, ".if ", "", 1)
	cmd = strings.TrimSpace(cmd)
	//cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	orStr := "||"
	andStr := "&&"
	eqStr := "=="
	nEqStr := "!="

	for {
		condVal1, condOp1 := getIfData(cmd, orStr, andStr, eqStr, nEqStr)
		condVal1 = strings.TrimSpace(condVal1)
		tVal1 = condVal1
		tOp1 = condOp1

		//fmt.Println(condVal1 + "," + condOp1)
		if condOp1 == "" && prevOp1 == "" {
			if condVal1 == "true" {
				return false, ""
			} else {
				return true, ""
			}
		}

		if prevOp1 != "" {
			switch prevOp1 {
			case orStr:
				if prevVal1 == "true" || condVal1 == "true" {
					condVal1 = "true"
				} else {
					condVal1 = "false"
				}
			case andStr:
				if prevVal1 == "true" && condVal1 == "true" {
					condVal1 = "true"
				} else {
					condVal1 = "false"
				}
			case eqStr:
				if prevVal1 == condVal1 {
					condVal1 = "true"
				} else {
					condVal1 = "false"
				}
			case nEqStr:
				if prevVal1 != condVal1 {
					condVal1 = "true"
				} else {
					condVal1 = "false"
				}
			}
		}

		if condOp1 == "" {
			if condVal1 == "true" {
				return false, ""
			} else {
				return true, ""
			}
		}

		cmd = strings.Replace(cmd, tVal1, "", 1)
		cmd = strings.TrimSpace(cmd)
		cmd = strings.Replace(cmd, tOp1, "", 1)
		cmd = strings.TrimSpace(cmd)

		prevVal1 = condVal1
		prevOp1 = condOp1

	}

}

func sorlOrchTag(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	//fmt.Println("inside...tag")
	cmd = strings.Replace(cmd, ".tag ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Replace(cmd, "{", "", 1)
	cmd = strings.TrimSpace(cmd)

	tags := (*allProp)["sr:tags"]

	if tags == "" {
		return false, cmd
	}

	for _, lCmd := range strings.Split(cmd, ",") {
		if strings.Contains(tags+",", lCmd+",") {
			return false, cmd
		}
	}

	return true, cmd
}

func sorlOrchPrint(cmd string, color string) {
	cmd = strings.Replace(cmd, ".print", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	sshPrint(color, cmd+"\n")
}

func sorlOrchSleep(cmd string) {

	cmd = strings.Replace(cmd, ".sleep", "", 1)
	cmd = strings.TrimSpace(cmd)
	lVal, err := strconv.Atoi(cmd)

	if err != nil {
		lVal = 1
	}
	time.Sleep(time.Second * time.Duration(lVal))
}

func sorlOrchVar(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {
	cmd = strings.Replace(cmd, ".var ", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	vars := strings.Split(cmd, "=")
	(*allProp)[vars[0]] = vars[1]

	//printMap("Var Map", SorlMap(allProp))
}

func sorlOrchWait(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (int, string) {

	cmd = strings.Replace(cmd, ".wait", "", 1)
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
