package main

import (
	"strings"
)

func procSorlOrchVar(ss *SorlSSH, cmds string, allProp *Property) {

	cmdsList := strings.Split(cmds, "\n")[1:]
	//PrintList("Var CMDS", cmdsList)
	cmds = ""
	for _, s := range cmdsList {
		cmds += ".var " + s + "\n"
	}

	ss.sorlOrchestration(cmds, allProp)

}

func procSorlOrchDebug(ss *SorlSSH, cmds string, allProp *Property) {

	cliOptDebug := false
	if (*allProp)["sr:debug"] == "true" {
		cliOptDebug = true
	}

	cmdsList := strings.Split(cmds, "\n")[1:]
	//PrintList("Debug CMDS", cmdsList)
	cmds = strings.Join(cmdsList, "\n")
	cmds = strings.TrimRight(cmds, "\n")

	if cliOptDebug {
		ss.sorlOrchestration(cmds, allProp)
	}

}

func procSorlOrchIf(ss *SorlSSH, cmds string, allProp *Property) {

	prevVal1 := ""
	prevOp1 := ""
	tVal1 := ""
	tOp1 := ""
	//
	cmdsList := strings.Split(cmds, "\n")[1:]
	//PrintList("IF CMDS", cmdsList)

	cmd := strings.Split(cmds, "\n")[0]
	cmd = strings.Replace(cmd, ".if ", "", 1)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.TrimRight(cmd, "{")
	cmd = strings.TrimSpace(cmd)

	cmds = strings.Join(cmdsList, "\n")
	cmds = strings.TrimRight(cmds, "\n")

	//fmt.Println("inside...tag")

	orStr := "||"
	andStr := "&&"
	eqStr := "=="
	nEqStr := "!="

	cmd, _ = replaceProp(cmd, (*allProp))
	//fmt.Println("IF Condition: " + cmd)

	for {
		condVal1, condOp1 := getIfData(cmd, orStr, andStr, eqStr, nEqStr)
		condVal1 = strings.TrimSpace(condVal1)
		tVal1 = condVal1
		tOp1 = condOp1

		//fmt.Println(condVal1 + "," + condOp1)
		if condOp1 == "" && prevOp1 == "" {
			if condVal1 == "true" {
				ss.sorlOrchestration(cmds, allProp)
				return
			} else {
				return
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
				ss.sorlOrchestration(cmds, allProp)
				return
			} else {
				return
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
