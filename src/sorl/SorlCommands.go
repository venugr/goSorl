package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func callSorlOrchUpper(cmd string, allProp *Property) {

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

func callSorlOrchLower(cmd string, allProp *Property) {

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

func callSorlOrchPrint(cmd string, allProp *Property) {
	color := (*allProp)["sr:color"]
	cmd = strings.Replace(cmd, ".println", "", 1)
	cmd = strings.Replace(cmd, ".print", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmd, _ = replaceProp(cmd, (*allProp))
	sshPrint(color, cmd)
}

func callSorlOrchPrintln(cmd string, allProp *Property) {
	callSorlOrchPrint(cmd+"\n", allProp)
}

func callSorlOrchVar(cmd string, allProp *Property) {

	//fmt.Println(cmd)

	if strings.HasPrefix(cmd, ".var ") && (!strings.HasSuffix(cmd, "{")) {

		cmd = strings.Replace(cmd, ".var ", "", 1)
		cmd = strings.TrimLeft(cmd, " ")
		cmd = strings.TrimLeft(cmd, "\t")
		cmd = strings.TrimLeft(cmd, " ")

		vars := strings.Split(cmd, "=")
		(*allProp)[vars[0]] = vars[1]

	}

	if strings.HasPrefix(cmd, ".var ") && strings.HasSuffix(cmd, "{") {

		//fmt.Println("===>" + cmd)
		(*allProp)["_block.started"] += ".var.yes,"
		(*allProp)["_block.names"] += ".var,"
		(*allProp)["_block.current"] = ".var"
	}

}

func callSorlOrchDebug(cmd string, allProp *Property) {

	(*allProp)["_block.started"] += ".debug.yes,"
	(*allProp)["_block.names"] += ".debug,"
	(*allProp)["_block.current"] = ".debug"

}

func callSorlOrchIf(cmd string, allProp *Property) {

	(*allProp)["_block.started"] += ".if.yes,"
	(*allProp)["_block.names"] += ".if,"
	(*allProp)["_block.current"] = ".if"

}

func callSorlOrchSleep(cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".sleep", "", 1)
	cmd = strings.TrimSpace(cmd)
	lVal, err := strconv.Atoi(cmd)

	if err != nil {
		lVal = 1
	}
	time.Sleep(time.Second * time.Duration(lVal))
}

func callSorlOrchShow(cmd string, allProp *Property) {

	cmd = strings.Replace(cmd, ".show", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimLeft(cmd, "{")
	cmd = strings.TrimRight(cmd, "}")
	cmd = strings.TrimSpace(cmd)
	sshPrint((*allProp)["sr:color"], "\n"+(*allProp)[cmd])

}

func callSorlOrchInput(cmd string, allProp *Property) {

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

func callSorlOrchName(cmd string, allProp *Property) {

	color := (*allProp)["sr:color"]
	cmd, _ = replaceProp(cmd, (*allProp))

	cmd = strings.Replace(cmd, ".name", "", 1)
	cmd = strings.TrimLeft(cmd, " ")
	cmdLen := len(cmd)

	sshPrint(color, "\n\n\n"+strings.Repeat("*", cmdLen+4)+"\n")
	sshPrint(color, "* "+cmd+" *\n")
	sshPrint(color, strings.Repeat("*", cmdLen+4)+"\n")

}

func callSorlOrchIncr(cmd string, allProp *Property) {
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

func callSorlOrchEcho(cmd string, allProp *Property) {
	cmd = strings.Replace(cmd, ".echo", "", 1)
	cmd = strings.TrimSpace(cmd)

	if cmd == "off" {
		(*allProp)["sr:echo"] = "off"
		//return false
	}

	(*allProp)["sr:echo"] = "on"
	//return true

}

func callSorlOrchDecr(cmd string, allProp *Property) {
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

func callSorlOrchFail(cmd string, allProp *Property) {

	tempCmdOut := (*allProp)["_cmd.temp.out"]
	cmd = strings.Replace(cmd, ".fail ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		(*allProp)["_fail.test"] = "true"
		return
	}

	(*allProp)["_fail.test"] = "false"

}

func callSorlOrchPass(cmd string, allProp *Property) {

	tempCmdOut := (*allProp)["_cmd.temp.out"]
	cmd = strings.Replace(cmd, ".pass ", "", 1)
	cmd = strings.TrimSpace(cmd)

	if strings.Contains(tempCmdOut, cmd) {
		(*allProp)["_pass.test"] = "true"
		return
	}

	(*allProp)["_pass.test"] = "false"

}

func callSorlOrchTest(cmd string, allProp *Property) {

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
