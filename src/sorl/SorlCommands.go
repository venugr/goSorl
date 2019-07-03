package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func callSorlOrchAnimate(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".animate", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	aniTimeStr := strings.Split(cmd, " ")[0]
	cmd = strings.TrimLeft(cmd, aniTimeStr)
	cmd = strings.TrimSpace(cmd)

	aniTime, _ := strconv.Atoi(aniTimeStr)

	for _, i := range cmd {
		fmt.Print(string(i))
		time.Sleep(time.Millisecond * time.Duration(aniTime))
	}

}

func callSorlOrchTrimleft(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".trimleft", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	propName1 := strings.Split(cmd, " ")[0]
	cmd = strings.TrimLeft(cmd, propName1)
	cmd = strings.TrimSpace(cmd)
	cmdStr, _ := replaceProp(cmd, (*allProp))

	(*allProp)[propName1] = strings.Join(strings.Split(cmdStr, "\n")[1:], "\n")

}

func callSorlOrchTrimright(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".trimright", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	propName1 := strings.Split(cmd, " ")[0]
	cmd = strings.TrimLeft(cmd, propName1)
	cmd = strings.TrimSpace(cmd)
	cmdStr, _ := replaceProp(cmd, (*allProp))
	iLen := len(strings.Split(cmdStr, "\n"))
	(*allProp)[propName1] = strings.Join(strings.Split(cmdStr, "\n")[:iLen-1], "\n")

}
func callSorlOrchSelect(ss *SorlSSH, cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".select", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	idx := strings.Index(cmd, "from")
	selStr := cmd[:idx]
	cmd = strings.TrimLeft(cmd, selStr)
	cmd = strings.TrimLeft(cmd, " ")
	cmd = strings.TrimLeft(cmd, "from")
	cmd = strings.TrimSpace(cmd)
	cmdStr, _ := replaceProp(cmd, (*allProp))

	selVals := strings.Split(selStr, ",")
	idxVals := ""
	selResStr := ""

	lCama := ""
	for _, lVal := range selVals {
		lVal = strings.TrimSpace(lVal)
		idxVals += lCama + lVal
		lCama = ","
	}

	var rExp = regexp.MustCompile(`\s+`)

	for _, lVal := range strings.Split(cmdStr, "\n") {

		if strings.TrimSpace(lVal) == "" {
			continue
		}

		lVal = string(rExp.ReplaceAll([]byte(lVal), []byte(" ")))
		words := strings.Split(lVal, " ")

		lLen := len(words)

		lSpace := ""

		for _, lKey := range strings.Split(idxVals, ",") {

			iKey, _ := strconv.Atoi(lKey)
			if (iKey - 1) < lLen {

				//fmt.Println("Idx:", lKey-1)
				selResStr += lSpace + words[iKey-1]
				lSpace = " "
			} else {
				selResStr += lSpace + "(*NA)"
				lSpace = " "

			}
		}

		if lSpace == " " {
			selResStr += "\n"
		}
	}

	(*allProp)["_select.result.str"] = selResStr

}

func callSorlOrchStyle(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".style", "", 1)
	cmd = strings.TrimLeft(cmd, " ")

	vars := strings.Split(cmd, " ")
	cmd = strings.TrimLeft(cmd, vars[0])
	cmd = strings.TrimLeft(cmd, " ")
	clrCode := ClrUnColor

	switch vars[0] {
	case "red":
		clrCode = ClrRed
	case "yellow":
		clrCode = ClrYellow
	case "green":
		clrCode = ClrGreen
	case "blue":
		clrCode = ClrBlue
	case "white":
		clrCode = ClrWhite
	case "magenta":
		clrCode = ClrMagenta
	case "cyan":
		clrCode = ClrCyan
	}

	cmd, _ = replaceProp(cmd, (*allProp))
	sshPrint(clrCode, cmd+"\n")

}

func callSorlOrchUpper(ss *SorlSSH, cmd string, allProp *Property) {

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

func callSorlOrchLower(ss *SorlSSH, cmd string, allProp *Property) {

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

func callSorlOrchPrint(ss *SorlSSH, cmd string, allProp *Property) {
	color := (*allProp)["sr:color"]
	cmd = strings.Replace(cmd, ".println", "", 1)
	cmd = strings.Replace(cmd, ".print", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmd, _ = replaceProp(cmd, (*allProp))

	sshPrint(color, cmd)
}

func callSorlOrchPrintln(ss *SorlSSH, cmd string, allProp *Property) {
	callSorlOrchPrint(ss, cmd+"\n", allProp)
}

func callSorlOrchVar(ss *SorlSSH, cmd string, allProp *Property) {

	//fmt.Println(cmd)

	if strings.HasPrefix(cmd, ".var ") && (!strings.HasSuffix(cmd, "{")) {

		cmd = strings.Replace(cmd, ".var ", "", 1)
		cmd = strings.TrimLeft(cmd, " ")
		cmd = strings.TrimLeft(cmd, "\t")
		cmd = strings.TrimLeft(cmd, " ")

		vars := strings.Split(cmd, "=")
		cmd = strings.TrimLeft(cmd, vars[0])
		cmd = strings.TrimLeft(cmd, " ")
		cmd = strings.TrimLeft(cmd, "=")
		(*allProp)[vars[0]] = cmd

	}

	if strings.HasPrefix(cmd, ".var ") && strings.HasSuffix(cmd, "{") {

		//fmt.Println("===>" + cmd)
		(*allProp)["_block.started"] += ".var.yes,"
		(*allProp)["_block.names"] += ".var,"
		(*allProp)["_block.current"] = ".var"
	}

}

func callSorlOrchDebug(ss *SorlSSH, cmd string, allProp *Property) {

	(*allProp)["_block.started"] += ".debug.yes,"
	(*allProp)["_block.names"] += ".debug,"
	(*allProp)["_block.current"] = ".debug"

}

func callSorlOrchIf(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".if ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_block.started"] += ".if.yes,"
	(*allProp)["_block.names"] += ".if,"
	(*allProp)["_block.current"] = ".if"
	(*allProp)["_block.current.if"] = cmd

}

func callSorlOrchFunc(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".func ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Replace(cmd, "{", "", 1)
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_block.started"] += ".func.yes,"
	(*allProp)["_block.names"] += ".func,"
	(*allProp)["_block.current"] = ".func"
	(*allProp)["_block.current.funcName"] = cmd
}

func callSorlOrchSleep(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".sleep", "", 1)
	cmd = strings.TrimSpace(cmd)
	lVal, err := strconv.Atoi(cmd)

	if err != nil {
		lVal = 1
	}
	time.Sleep(time.Second * time.Duration(lVal))
}

func callSorlOrchShow(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".show", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimLeft(cmd, "{")
	cmd = strings.TrimRight(cmd, "}")
	cmd = strings.TrimSpace(cmd)
	sshPrint((*allProp)["sr:color"], "\n"+(*allProp)[cmd])

}

func callSorlOrchInput(ss *SorlSSH, cmd string, allProp *Property) {

	color := (*allProp)["sr:color"]
	tCmd := cmd
	cmd = strings.Replace(cmd, ".input", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	lPropList := strings.Split(cmd, " ")

	if len(lPropList) == 0 {
		fmt.Println(errors.New(".input command is ill formed: " + tCmd))
	}

	propName := lPropList[0]
	cmd = strings.Replace(cmd, propName, "", 1)
	cmd = strings.TrimLeft(cmd, " ")

	reader := bufio.NewReader(os.Stdin)
	sshPrint(color, cmd+" ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")

	(*allProp)[propName] = text

}

func callSorlOrchName(ss *SorlSSH, cmd string, allProp *Property) {

	color := (*allProp)["sr:color"]
	cmd, _ = replaceProp(cmd, (*allProp))

	cmd = strings.Replace(cmd, ".name", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmdLen := len(cmd)

	sshPrint(color, "\n\n\n"+strings.Repeat("*", cmdLen+4)+"\n")
	sshPrint(color, "* "+cmd+" *\n")
	sshPrint(color, strings.Repeat("*", cmdLen+4)+"\n")

}

func callSorlOrchIncr(ss *SorlSSH, cmd string, allProp *Property) {
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

func callSorlOrchEcho(ss *SorlSSH, cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".echo", "", 1)
	cmd = strings.TrimSpace(cmd)

	if cmd == "off" {
		(*allProp)["sr:echo"] = "off"
		//return false
	}

	(*allProp)["sr:echo"] = "on"
	//return true

}

func callSorlOrchDecr(ss *SorlSSH, cmd string, allProp *Property) {
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

func callSorlOrchFail(ss *SorlSSH, cmd string, allProp *Property) {

	tempCmdOut := (*allProp)["_cmd.temp.out"]
	cmd = strings.Replace(cmd, ".fail ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		(*allProp)["_fail.test"] = "true"
		return
	}

	(*allProp)["_fail.test"] = "false"

}

func callSorlOrchPass(ss *SorlSSH, cmd string, allProp *Property) {

	tempCmdOut := (*allProp)["_cmd.temp.out"]
	cmd = strings.Replace(cmd, ".pass ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		(*allProp)["_pass.test"] = "true"
		return
	}

	(*allProp)["_pass.test"] = "false"

}

func callSorlOrchTest(ss *SorlSSH, cmd string, allProp *Property) {

	tempCmdOut := (*allProp)["_cmd.temp.out"]

	cmd = strings.Replace(cmd, ".test ", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	propName := strings.Split(cmd, " ")[0]

	cmd = strings.TrimLeft(cmd, propName)
	cmd = strings.TrimLeft(cmd, " ")

	(*allProp)[propName] = "false"
	if strings.Contains(tempCmdOut, cmd) {
		(*allProp)[propName] = "true"

	}

}

func callSorlOrchCall(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".call", "", 1)
	cmd = strings.TrimSpace(cmd)
	ss.sorlOrchestration((*allProp)["_func.name."+cmd], allProp)

}

func callSorlOrchTag(ss *SorlSSH, cmd string, allProp *Property) {

	//fmt.Println("inside...tag")
	cmd = strings.Replace(cmd, ".tag ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Replace(cmd, "{", "", 1)
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_block.started"] += ".tag.yes,"
	(*allProp)["_block.names"] += ".tag,"
	(*allProp)["_block.current"] = ".tag"
	(*allProp)["_block.current.tag"] = cmd
}

func callSorlOrchRange(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".range ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_block.started"] += ".range.yes,"
	(*allProp)["_block.names"] += ".range,"
	(*allProp)["_block.current"] = ".range"
	(*allProp)["_block.current.range"] = cmd
}

func callSorlOrchWhile(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".while ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_block.started"] += ".while.yes,"
	(*allProp)["_block.names"] += ".while,"
	(*allProp)["_block.current"] = ".while"
	(*allProp)["_block.current.while"] = cmd

}

func callSorlOrchLoad(ss *SorlSSH, cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".load ", "", 1)
	cmd = strings.TrimSpace(cmd)

	loadFile, _ := replaceProp(cmd, (*allProp))

	(*allProp)["sr:loadfile"] = loadFile
	(*allProp)["sr:load"] = "yes"
	ss.sorlRunOrchestration(allProp)
}

func callSorlOrchReturn(ss *SorlSSH, cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".return", "", 1)
	cmd = strings.TrimSpace(cmd)

	(*allProp)["_return"] = "true"
	(*allProp)["_return.code"] = cmd

}

func callSorlOrchEnter(ss *SorlSSH, cmd string, allProp *Property) {

}

func callSorlOrchReplace(ss *SorlSSH, cmd string, allProp *Property) {

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
