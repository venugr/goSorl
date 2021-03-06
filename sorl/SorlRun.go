package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func sorlRunOrchestration(session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {

	//fmt.Println("Run Orchestration....")
	(*allProp)["_cmd.output"] = ""

	orchFile := string((*allProp)["sr:orchfile"])
	loadFile := string((*allProp)["sr:loadfile"])
	loadOk := string((*allProp)["sr:load"])

	if loadOk == "yes" {
		orchFile = loadFile
	}

	if _, ok := (*allProp)["_wait.prev.cmd"]; !ok {
		(*allProp)["_wait.prev.cmd"] = ""
	}

	if _, ok := (*allProp)["_wait.done"]; !ok {
		(*allProp)["_wait.done"] = "-1"
	}

	(*allProp)["_tags.order"] = ""
	(*allProp)["_last.tag"] = ""

	commands, _ := ReadFile(orchFile)

	//fmt.Println("==>1." + (*allProp)["sr:debug"] + "<==")

	sorlOrchestration(strings.Join(commands, "\n"), session, sshIn, sshOut, allProp)
}

func sorlOrchestration(cmdLines string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) {

	//fmt.Println("Run Orchestration....")
	(*allProp)["_cmd.output"] = ""

	//orchFile := string((*allProp)["sr:orchfile"])
	//color := string((*allProp)["sr:color"])
	keepNoCmdLogs, _ := strconv.Atoi((*allProp)["sr:keep"])
	keepCmdLogs := make([]string, keepNoCmdLogs)
	//loadFile := string((*allProp)["sr:loadfile"])
	loadOk := string((*allProp)["sr:load"])
	//display := string(allProp["sr:display"])

	skipTagLines := false
	skipIfLines := false
	skipDebugLines := false
	//tagName := ""

	//if loadOk == "yes" {
	//	orchFile = loadFile
	//}

	//rangeSeq := 0

	commands := strings.Split(cmdLines, "\n")

	cmdOut := ""
	prevWaitCmd := (*allProp)["_wait.prev.cmd"]
	runWaitOk := false
	tempCmdOut := ""
	skipVarLines := false
	skipFuncLines := false
	skipRangeLines := false
	skipWhileLines := false
	whileCondStr := ""
	rangePropName := ""
	funcName := ""
	funcLoops := 0
	whileLoops := 0
	isRemoved := false
	waitMatchId := -1
	ifReq := false
	tagsOrder := (*allProp)["_tags.order"]
	lastTag := (*allProp)["_last.tag"]
	waitDone := (*allProp)["_wait.done"]
	//echoOn := true
	(*allProp)["sr:echo"] = "on"

	//fmt.Println("==>2." + (*allProp)["sr:debug"] + "<==")

	mapFuncs := map[string]string{}
	mapRanges := map[string]string{}
	mapWhile := map[string]string{}

	if loadOk == "no" {
		//waitFor(color, []string{"$", "[BAN83] ?"}, sshIn)
	}

	for _, cmd := range commands {

		//fmt.Println("Cmd:", cmd)

		cmd = strings.TrimLeft(cmd, " ")

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

		isTag := strings.HasPrefix(cmd, "}") || strings.HasPrefix(cmd, ".if") || strings.HasPrefix(cmd, ".func")
		isTag = isTag || strings.HasPrefix(cmd, ".var") || strings.HasPrefix(cmd, ".tag") || strings.HasPrefix(cmd, ".range")
		isTag = isTag || strings.HasPrefix(cmd, ".while")

		if (!isTag) && (skipTagLines || skipIfLines) {
			//fmt.Println("Skipping...:", cmd)
			continue
		}

		if strings.HasPrefix(cmd, "}") {
			//fmt.Println("Skipping...1:", cmd)
			//fmt.Println("1.", tagsOrder, lastTag)

			if whileLoops == 0 {
				tagsOrder = strings.TrimSuffix(tagsOrder, lastTag)
				tagsOrder = strings.TrimSuffix(tagsOrder, "skip.")
				(*allProp)["_tags.order"] = tagsOrder
			}

			if skipDebugLines && strings.EqualFold(lastTag, "debug,") && (!strings.Contains(tagsOrder, "skip.debug,")) {
				skipDebugLines = false

			}

			if skipTagLines && strings.EqualFold(lastTag, "tag,") && (!strings.Contains(tagsOrder, "skip.tag,")) {
				skipTagLines = false

			}

			if skipIfLines && strings.EqualFold(lastTag, "if,") && (!strings.Contains(tagsOrder, "skip.if,")) {
				skipIfLines = false
			}

			if skipVarLines && strings.EqualFold(lastTag, "var,") {
				skipVarLines = false
			}

			if skipFuncLines && strings.EqualFold(lastTag, "func,") {
				//skipFuncLines = false
			}

			if skipRangeLines && strings.EqualFold(lastTag, "range,") {
				//skipFuncLines = false
				rangeValue, _ := replaceProp(rangePropName, Property(*allProp))
				for _, idxVal := range strings.Split(rangeValue, "\n") {
					idxVal = strings.TrimRight(idxVal, string(10))
					idxVal = strings.TrimRight(idxVal, string(13))

					(*allProp)["range.value"] = idxVal
					sorlOrchestration((*allProp)["_range."+rangePropName], session, sshIn, sshOut, allProp)
				}

				ifReq = true
				skipRangeLines = false

			}

			if whileLoops == 0 {
				if tagsOrder != "" {
					tagsList := strings.Split(tagsOrder, ",")
					lastTag = tagsList[len(tagsList)-2] + ","
					lastTag = strings.TrimPrefix(lastTag, "skip.")
				} else {
					lastTag = ""
				}

				(*allProp)["_tags.order"] = tagsOrder
				(*allProp)["_last.tag"] = lastTag
				//fmt.Println("2.", tagsOrder, lastTag)
				//fmt.Println("1.True/False", skipTagLines, skipIfLines)
			}
			/*
				if !skipFuncLines {
					continue
				}

				if !skipRangeLines {
					continue
				}

				if !skipWhileLines {
					continue
				}
			*/
			if !(skipFuncLines || skipRangeLines || skipWhileLines) {
				continue
			}

			//fmt.Println("I AM HERE...")
		}

		//fmt.Println("2.True/False", skipTagLines, skipIfLines, skipFuncLines)

		if skipVarLines && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchVar(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		if skipFuncLines && strings.HasSuffix(strings.TrimRight(cmd, " "), "{") {
			funcLoops++
		}

		if skipWhileLines && strings.HasSuffix(strings.TrimRight(cmd, " "), "{") {
			whileLoops++
		}

		bracStr := strings.TrimSpace(cmd)
		//fmt.Printf("Cmd:==>%v<==", bracStr)
		if skipFuncLines && (bracStr == "}") {
			funcLoops--
			if funcLoops == 0 {
				skipFuncLines = false
				continue
			}

		}

		if skipWhileLines && (bracStr == "}") {
			whileLoops--
			if whileLoops == 0 {
				skipWhileLines = false

				for {
					whileValue, _ := replaceProp(whileCondStr, Property(*allProp))
					//fmt.Println("WhileValue: " + whileValue)
					lOk, _ := sorlOrchIf(whileValue, session, sshIn, sshOut, allProp)
					//fmt.Println("WhileOk: ", lOk)

					if lOk {
						break
					}

					//fmt.Println((*allProp)["_while."+whileCondStr])
					//fmt.Println("2.", tagsOrder, lastTag)
					sorlOrchestration((*allProp)["_while."+whileCondStr], session, sshIn, sshOut, allProp)
				}
				ifReq = true
				skipWhileLines = false

				continue
			}

		}

		//fmt.Println("Ture/Value:", skipFuncLines, funcLoops)

		if skipFuncLines && funcLoops != 0 {
			mapFuncs[funcName] += cmd + "\n"
			(*allProp)["_func."+funcName] = mapFuncs[funcName]
			continue
		}

		if skipRangeLines {
			mapRanges["_range."+rangePropName] += cmd + "\n"
			(*allProp)["_range."+rangePropName] = mapRanges["_range."+rangePropName]
			continue
		}

		if skipWhileLines {
			mapWhile["_while."+whileCondStr] += cmd + "\n"
			(*allProp)["_while."+whileCondStr] = mapWhile["_while."+whileCondStr]
			//fmt.Println("Add/While:" + cmd)
			//fmt.Println("While Loop: ", whileLoops)
			continue
		}

		if runWaitOk && (!strings.HasPrefix(cmd, ".wait")) {
			waitMatchId, cmdOut = sorlOrchWait(prevWaitCmd, session, sshIn, sshOut, allProp)
			(*allProp)["_wait.match.id"] = strconv.Itoa(waitMatchId)
			tempCmdOut += cmdOut
			(*allProp)["_cmd.output"] = tempCmdOut
			cmdList := strings.Split(cmdOut, "\n")
			cmdListLen := len(cmdList) - 1
			(*allProp)["_wait.matched.prompt"] = cmdList[cmdListLen]
			//waitDone = "0"
			//(*allProp)["_wait.done"] = "0"

			/*
				tEchoOn := (*allProp)["sr:echo"]
				(*allProp)["sr:echo"] = "off"
				runShellCmd("echo $?", sshOut)
				sorlOrchWait(prevWaitCmd, session, sshIn, sshOut, allProp)
				(*allProp)["sr:echo"] = tEchoOn
			*/

		}
		runWaitOk = false

		/*
			if !isTag {
				continue
			}
		*/

		isRemoved = false
		if strings.HasSuffix(cmd, "{") {
			cmd = strings.TrimRight(cmd, "{")
			isRemoved = true
		}

		if strings.HasPrefix(cmd, ".show ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			lProp := sorlOrchShow(cmd)
			sshPrint((*allProp)["sr:color"], "\n"+(*allProp)[lProp])
			ifReq = true
			continue
		}

		if strings.HasPrefix(cmd, ".setwait ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			prevWaitCmd = strings.Replace(cmd, ".setwait", ".wait", 1)
			(*allProp)["_wait.prev.cmd"] = prevWaitCmd
			(*allProp)["_wait.string"] = strings.TrimSpace(strings.Replace(cmd, ".setwait ", "", 1))
			continue
		}

		if strings.HasPrefix(cmd, ".upper ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchUpper(cmd, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".lower ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchLower(cmd, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".replace ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchReplace(cmd, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".incr ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchIncr(cmd, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".decr ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchDecr(cmd, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".echo ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			//ifReq = false
			if sorlOrchEcho(cmd, allProp) {
				ifReq = true
			}
			continue
		}

		if strings.HasPrefix(cmd, ".while ") {
			if !skipTagLines {
				skipWhileLines, whileCondStr = sorlOrchWhile(cmd, session, sshIn, sshOut, allProp)
				//rangeSeq += 1
				mapWhile["_while."+whileCondStr] = ""
				(*allProp)["_while."+whileCondStr] = ""
			}
			tagsOrder += "while,"
			lastTag = "while,"

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			whileLoops++
			//funcLoops++
			continue
		}

		err1 := errors.New("")
		if !(skipTagLines || skipIfLines || skipDebugLines) {
			cmd, err1 = replaceProp(cmd, Property(*allProp))
			checkError(err1)
		}

		if isRemoved {
			cmd = cmd + "{"
		}

		if strings.HasPrefix(cmd, ".name") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchName(cmd, (*allProp)["sr:color"])
			ifReq = true
			continue
		}

		if strings.HasPrefix(cmd, ".return") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			retCode := sorlOrchReturn(cmd)

			if retCode != 0 {
				fmt.Printf("return: %v\n", retCode)
				session.Close()
			}
			return
		}

		if strings.HasPrefix(cmd, ".sleep") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchSleep(cmd)
			continue
		}

		if strings.HasPrefix(cmd, ".clear") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			tempCmdOut = ""
			(*allProp)["_cmd.output"] = tempCmdOut
			continue
		}

		if strings.HasPrefix(cmd, ".println") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			//fmt.Printf("Cmd: --->" + cmd + "<----\n")
			sorlOrchPrintln(cmd, (*allProp)["sr:color"])
			ifReq = true
			continue
		}

		if strings.HasPrefix(cmd, ".print") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchPrint(cmd, (*allProp)["sr:color"])
			ifReq = true
			continue
		}

		if strings.HasPrefix(cmd, ".input ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchInput(cmd, (*allProp)["sr:color"], allProp)
			ifReq = true
			continue
		}

		if strings.HasPrefix(cmd, ".call ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			callName, _ := sorlOrchCall(cmd, (*allProp)["sr:color"], allProp)
			//ifReq = true

			//fmt.Println((*allProp)["_func."+callName])
			sorlOrchestration((*allProp)["_func."+callName], session, sshIn, sshOut, allProp)
			continue
		}

		if strings.HasPrefix(cmd, ".pass") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			if sorlOrchPass(cmd, (*allProp)["sr:color"], tempCmdOut) {
				sshPrint((*allProp)["sr:color"], "\n"+cmd+" : Failed\n")
				ifReq = true
				session.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(cmd, ".fail") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			if sorlOrchFail(cmd, (*allProp)["sr:color"], tempCmdOut) {
				sshPrint((*allProp)["sr:color"], "\n"+cmd+" : Failed\n")
				ifReq = true
				session.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(cmd, ".test") && (!(skipTagLines || skipIfLines || skipDebugLines)) {

			propName, testOk := sorlOrchTest(cmd, (*allProp)["sr:color"], tempCmdOut)

			(*allProp)[propName] = "false"

			if testOk {
				(*allProp)[propName] = "true"
			}

			continue
		}

		//prevWaitCmd = cmd

		if strings.HasPrefix(cmd, ".range ") {
			if !skipTagLines {
				skipRangeLines, rangePropName = sorlOrchRange(cmd, session, sshIn, sshOut, allProp)
				//rangeSeq += 1
				mapRanges["_range."+rangePropName] = ""
				(*allProp)["_range."+rangePropName] = ""
			}
			tagsOrder += "range,"
			lastTag = "range,"
			//funcLoops++

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".func ") {
			if !skipTagLines {
				skipFuncLines, funcName = sorlOrchFunc(cmd, session, sshIn, sshOut, allProp)
				mapFuncs[funcName] = ""
				(*allProp)["_func."+funcName] = ""
			}
			tagsOrder += "func,"
			lastTag = "func,"
			funcLoops++

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".tag ") {
			if !skipTagLines {
				skipTagLines, _ = sorlOrchTag(cmd, session, sshIn, sshOut, allProp)
				if skipTagLines {
					tagsOrder += "skip."
				}
			}
			tagsOrder += "tag,"
			lastTag = "tag,"

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".debug ") {
			//fmt.Println("==>3." + (*allProp)["sr:debug"] + "<==")
			if !skipTagLines {

				//fmt.Println("Debug block found!")
				//time.Sleep(1 * time.Second)
				skipDebugLines, _ = sorlOrchDebug(cmd, session, sshIn, sshOut, allProp)
				if skipDebugLines {
					//fmt.Println("Debug no run!")
					//time.Sleep(1 * time.Second)
					tagsOrder += "skip."
				}
			}
			tagsOrder += "debug,"
			lastTag = "debug,"
			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".if ") {
			if !skipTagLines {
				skipIfLines, _ = sorlOrchIf(cmd, session, sshIn, sshOut, allProp)
				if skipIfLines {
					tagsOrder += "skip."
				}
			}
			tagsOrder += "if,"
			lastTag = "if,"

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".var ") && (!strings.Contains(cmd, "{")) && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			//fmt.Printf("var single %v", cmd)

			sorlOrchVar(cmd, session, sshIn, sshOut, allProp)

			//printMap("Var Map", SorlMap(*allProp))
			continue
		}

		if strings.HasPrefix(cmd, ".var ") && strings.Contains(cmd, "{") {
			if !skipTagLines {
				skipVarLines = true
			}
			tagsOrder += "var,"
			lastTag = "var,"
			//fmt.Printf("var group")

			(*allProp)["_tags.order"] = tagsOrder
			(*allProp)["_last.tag"] = lastTag
			continue
		}

		if strings.HasPrefix(cmd, ".load ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			sorlOrchLoad(cmd, session, sshIn, sshOut, allProp)
			continue
		}

		cmdOut = ""
		if strings.HasPrefix(cmd, ".wait ") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			waitMatchId, cmdOut = sorlOrchWait(cmd, session, sshIn, sshOut, allProp)
			prevWaitCmd = cmd
			(*allProp)["_wait.prev.cmd"] = prevWaitCmd
			(*allProp)["_wait.string"] = strings.TrimSpace(strings.Replace(cmd, ".wait ", "", 1))
			(*allProp)["_wait.match.id"] = strconv.Itoa(waitMatchId)
			runWaitOk = false
			tempCmdOut += cmdOut
			(*allProp)["_cmd.output"] = tempCmdOut
			cmdList := strings.Split(cmdOut, "\n")
			cmdListLen := len(cmdList) - 1
			(*allProp)["_wait.matched.prompt"] = cmdList[cmdListLen]
			//waitDone = "0"
			//(*allProp)["_wait.done"] = "0"
			//ifReq = false

			/*
				tEchoOn := (*allProp)["sr:echo"]
				(*allProp)["sr:echo"] = "off"
				runShellCmd("echo $?", sshOut)
				sorlOrchWait(cmd, session, sshIn, sshOut, allProp)
				(*allProp)["sr:echo"] = tEchoOn
			*/

			continue
		}

		if strings.HasPrefix(cmd, ".enter") && (!(skipTagLines || skipIfLines || skipDebugLines)) {
			cmd = ""
		}

		if skipTagLines || skipIfLines || skipDebugLines {
			continue
		}

		runWaitOk = true
		//fmt.Println("R: Cmd:", cmd)

		//color := (*allProp)["sr:color"]
		//display := (*allProp)["sr:display"]
		//if ifReq && waitDone != "-1" {
		//	sshPrint((*allProp)["sr:color"], "\n"+(*allProp)["_wait.matched.prompt"])
		//}
		//ifReq = false
		waitDone = "0"
		(*allProp)["_wait.done"] = "0"
		//fmt.Println("R: Cmd:", cmd)
		lCmd := strings.ReplaceAll(cmd, " ", "")
		if strings.Contains(lCmd, "rm-rf*") {
			sshPrint((*allProp)["sr:color"], "\nsorl: can not process "+cmd)
			reader := bufio.NewReader(os.Stdin)
			sshPrint((*allProp)["sr:color"], "\nDo you want to proceed(yes/no)? ")
			yesNo, _ := reader.ReadString('\n')
			yesNo = strings.TrimRight(yesNo, "\n")
			ifReq = true
			if yesNo != "yes" {
				runWaitOk = false
				continue
			}
		}

		if ifReq && waitDone != "-1" {
			sshPrint((*allProp)["sr:color"], "\n"+(*allProp)["_wait.matched.prompt"])
		}
		ifReq = false

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

func sorlOrchTest(cmd, color, tempCmdOut string) (string, bool) {

	cmd = strings.Replace(cmd, ".test ", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	propName := strings.Split(cmd, " ")[0]

	cmd = strings.TrimLeft(cmd, propName)
	cmd = strings.TrimLeft(cmd, " ")

	if strings.Contains(tempCmdOut, cmd) {
		return propName, true
	}

	return propName, false

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

func sorlOrchDebug(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {
	cmd = strings.Replace(cmd, ".debug ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	//fmt.Println("==>" + (*allProp)["sr:debug"] + "<==")
	//time.Sleep(1 * time.Second)
	if (*allProp)["sr:debug"] == "true" {
		return false, cmd
	}

	return true, cmd

}

func sorlOrchIf(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	prevVal1 := ""
	prevOp1 := ""
	tVal1 := ""
	tOp1 := ""
	//fmt.Println("inside...tag")
	cmd = strings.Replace(cmd, ".if ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
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

func sorlOrchWhile(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	cmd = strings.Replace(cmd, ".while ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	return true, cmd
}

func sorlOrchRange(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	cmd = strings.Replace(cmd, ".range ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	return true, cmd
}

func sorlOrchFunc(cmd string, session *ssh.Session, sshIn io.Reader, sshOut io.WriteCloser, allProp *Property) (bool, string) {

	cmd = strings.Replace(cmd, ".func ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Replace(cmd, "{", "", 1)
	cmd = strings.TrimSpace(cmd)

	return true, cmd
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
		if strings.Contains(","+tags+",", ","+lCmd+",") {
			return false, cmd
		}
	}

	return true, cmd
}

func sorlOrchCall(cmd string, color string, allProp *Property) (string, error) {

	cmd = strings.Replace(cmd, ".call", "", 1)
	cmd = strings.TrimSpace(cmd)

	return cmd, nil

}

func sorlOrchInput(cmd string, color string, allProp *Property) error {

	tCmd := cmd
	cmd = strings.Replace(cmd, ".input", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	lPropList := strings.Split(cmd, " ")

	if len(lPropList) == 0 {
		return errors.New(".input command is ill formed: " + tCmd)
	}

	propName := lPropList[0]
	cmd = strings.Replace(cmd, propName, "", 1)
	cmd = strings.TrimLeft(cmd, " ")

	reader := bufio.NewReader(os.Stdin)
	sshPrint(color, cmd+" ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")

	(*allProp)[propName] = text

	return nil

}

func sorlOrchUpper(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".upper", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	idx := strings.Index(cmd, " ")
	propName := cmd[:idx+1]
	propName = strings.TrimSpace(propName)

	cmd = strings.Replace(cmd, propName, "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmd, _ = replaceProp(cmd, (*allProp))

	(*allProp)[propName] = strings.ToUpper(cmd)

}

func sorlOrchLower(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".lower", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	idx := strings.Index(cmd, " ")
	propName := cmd[:idx+1]
	propName = strings.TrimSpace(propName)

	cmd = strings.Replace(cmd, propName, "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmd, _ = replaceProp(cmd, (*allProp))

	(*allProp)[propName] = strings.ToLower(cmd)

}

func sorlOrchIncr(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".incr", "", 1)
	cmd = strings.TrimSpace(cmd)

	propName := strings.TrimLeft(cmd, "{")
	propName = strings.TrimRight(propName, "}")
	propName = strings.TrimSpace(propName)

	cmd, _ = replaceProp(cmd, (*allProp))

	lVal, _ := strconv.Atoi(cmd)

	lVal++

	(*allProp)[propName] = strconv.Itoa(lVal)

}

func sorlOrchEcho(cmd string, allProp *Property) bool {
	cmd = strings.Replace(cmd, ".echo", "", 1)
	cmd = strings.TrimSpace(cmd)

	if cmd == "off" {
		(*allProp)["sr:echo"] = "off"
		return false
	}

	(*allProp)["sr:echo"] = "on"
	return true

}

func sorlOrchDecr(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".decr", "", 1)
	cmd = strings.TrimSpace(cmd)

	propName := strings.TrimLeft(cmd, "{")
	propName = strings.TrimRight(propName, "}")
	propName = strings.TrimSpace(propName)

	cmd, _ = replaceProp(cmd, (*allProp))

	lVal, _ := strconv.Atoi(cmd)

	lVal--

	(*allProp)[propName] = strconv.Itoa(lVal)

}

func sorlOrchReplace(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".replace", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	idx := strings.Index(cmd, " ")
	propName := cmd[:idx+1]
	propName = strings.TrimSpace(propName)

	cmd = strings.Replace(cmd, propName, "", 1)
	cmd = strings.TrimLeft(cmd, " ")

	idx = strings.Index(cmd, " ")
	srcProp := cmd[:idx+1]
	srcProp = strings.TrimSpace(srcProp)

	cmd = strings.Replace(cmd, srcProp, "", 1)
	cmd = strings.TrimLeft(cmd, " ")

	idx = strings.Index(cmd, " ")
	oldProp := cmd[:idx+1]
	oldProp = strings.TrimSpace(oldProp)

	newProp := strings.Replace(cmd, oldProp, "", 1)
	newProp = strings.TrimSpace(newProp)

	srcProp, _ = replaceProp(srcProp, (*allProp))
	oldProp, _ = replaceProp(oldProp, (*allProp))
	newProp, _ = replaceProp(newProp, (*allProp))

	(*allProp)[propName] = strings.ReplaceAll(srcProp, oldProp, newProp)

}

func sorlOrchPrint(cmd string, color string) {
	cmd = strings.Replace(cmd, ".println", "", 1)
	cmd = strings.Replace(cmd, ".print", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	sshPrint(color, cmd)
}

func sorlOrchPrintln(cmd string, color string) {
	sorlOrchPrint(cmd+"\n", color)
}

func sorlOrchReturn(cmd string) int {
	cmd = strings.Replace(cmd, ".return", "", 1)
	cmd = strings.TrimSpace(cmd)
	retCode, err := strconv.Atoi(cmd)

	if err != nil {
		return -1
	}

	return retCode
}
func sorlOrchName(cmd string, color string) {

	cmd = strings.Replace(cmd, ".name", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmdLen := len(cmd)

	sshPrint(color, "\n\n\n"+strings.Repeat("*", cmdLen+4)+"\n")
	sshPrint(color, "* "+cmd+" *\n")
	sshPrint(color, strings.Repeat("*", cmdLen+4)+"\n")

}

func sorlOrchShow(cmd string) string {

	cmd = strings.Replace(cmd, ".show", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimLeft(cmd, "{")
	cmd = strings.TrimRight(cmd, "}")
	cmd = strings.TrimSpace(cmd)
	return cmd

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
	echoOn := false

	if (*allProp)["sr:echo"] == "on" {
		echoOn = true
	}

	return waitFor(echoOn, color, display, waitStr, sshIn)

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
