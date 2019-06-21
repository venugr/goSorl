package main

import (
	"fmt"
	"strings"
	"time"
)

func callTest1() {
	ss := &SorlSSH{}
	alp := Property{}

	ss.sorlOrchestration(".upper h.u helloupper12\n   .lower   h.l HELLOLOWER89\n  .upper h.hu \n.println {h.l}\n.println {h.u}", &alp)
	fmt.Println(alp["h.u"])
	fmt.Println(alp["h.l"])
	fmt.Println(alp["h.hu"])

}

func callTest2() {
	ss := &SorlSSH{}
	alp := Property{}

	ss.sorlOrchestration(".var {\n   a.a=Apple\n   b.bb=Bat\n     c.123.abc=1234 * as 345\n   }\n#This is  a comment\n.println b.bb={b.bb}", &alp)

	fmt.Println(alp["a.a"])
	fmt.Println(alp["b.bb"])
	fmt.Println(alp["c.123.abc"])

}

func callTest3() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for Debug\n.println It is Debug\n.debug  {\n .println How is the Debug\n.var {\n a.b=AB\nb.c=BC\n		c.z=C2Z\n  }\n"
	cmd += ".println bebug is done\n}\n.println \n.println c.z={c.z}\n.println Debug Over/out"
	ss.sorlOrchestration(cmd, &alp)

	fmt.Println(alp["c.z"])
}

func callTest4() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for If\n.println It is IF\n.var n.n=10\n.if {n.n} == 10  {\n .println How is the IF\n"
	cmd += ".println IF is done\n}\n.println n.n={n.n}\n.println\n"
	cmd += ".if {n.n} == 10 {\n  .println {n.n} matched.\n}\n.println IF DONE"
	ss.sorlOrchestration(cmd, &alp)

	fmt.Println(alp["c.z"])
}

func callTest5() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for If\n"
	cmd += ".println Nested IF\n"
	cmd += ".if 10 == 10 {\n"
	cmd += "   .println 10 == 10\n"
	cmd += "   .if 5 == 10 {\n"
	cmd += "      .println 5 == 10 \n"
	cmd += "   }\n"
	cmd += "   .println Middle of If\n"
	cmd += "   .if tt == tt {\n"
	cmd += "      .println TT == TT\n"
	cmd += "   }\n"
	cmd += "   .println IF 10 == 10 Done\n"
	cmd += "}\n"
	cmd += ".println\n"
	cmd += "ls -ltr\n"
	cmd += ".println Hello Linux\n"
	cmd += "df -h"

	ss.sorlOrchestration(cmd, &alp)

	fmt.Println(alp["c.z"])
}

func callTest6() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .sleep\n"
	cmd += ".println sleep test\n"
	cmd += ".println sleep for 5s\n"
	cmd += ".sleep 5\n"
	cmd += ".var helo.help=Hello World, Welcome to Go Lang...\n"
	cmd += ".println another wait sleep for 5s\n"
	cmd += ".sleep 5\n"
	cmd += ".show {helo.help}\n"
	cmd += ".println "

	ss.sorlOrchestration(cmd, &alp)

}

func callTest7() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .input\n"
	cmd += ".println input test\n"
	cmd += ".input sl.emp.no Enter Emp Number? \n"
	cmd += ".input sl.emp.name Enter Emp Name? \n"
	cmd += ".println\n"
	cmd += ".println Num={sl.emp.no}\n"
	cmd += ".println Name={sl.emp.name}\n"
	cmd += ".println\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest8() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .name\n"
	cmd += ".println name test\n"
	cmd += ".input sl.emp.no Enter Emp Number? \n"
	cmd += ".input sl.emp.name Enter Emp Name? \n"
	cmd += ".println\n"
	cmd += ".println Num={sl.emp.no}\n"
	cmd += ".println Name={sl.emp.name}\n"
	cmd += ".name {sl.emp.no}\n"
	cmd += ".name {sl.emp.name}\n"
	cmd += ".println\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest9() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .incr/.decr/.echo\n"
	cmd += ".println name .incr/.decr/.echo\n"
	cmd += ".input sl.no1 Enter Emp Number? \n"
	cmd += ".println\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".incr {sl.no1}\n.incr {sl.no1}\n"
	cmd += ".println Num={sl.no1}\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".decr {sl.no1}\n.decr {sl.no1}\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".println\n"
	cmd += ".echo off\n"
	cmd += ".println Echo off\n"
	cmd += ".echo on\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest10() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .pass/.fail\n"

	cmd += ".var  _cmd.temp.out=This is test logoutin\n"
	cmd += ".pass logoutin\n"
	cmd += ".name {_cmd.temp.out}\n"
	cmd += ".fail logouin\n"
	cmd += ".name NAME={_cmd.temp.out}\n"
	cmd += ".test h.t1 logoutin\n"
	cmd += ".name {h.t1}\n"
	cmd += ".test h.t2 logotin\n"
	cmd += ".name {h.t2}\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest11() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .func\n"

	cmd += ".println funcname first\n"
	cmd += "# start the func name\n"
	cmd += ".func firstname {\n"
	cmd += "   .println Hello, I as in func 'firstname'\n"
	cmd += "   .if 10 != 10 {\n"
	cmd += "      .println 10 is not 1\n"
	cmd += "      .println good..and done.\n"
	cmd += "   }\n"
	cmd += "   .println After If inside func\n"
	cmd += "}\n"
	cmd += ".println Done Function..\n"
	cmd += ".name  FuncName: firstname\n"
	cmd += ".show {_func.name.firstname}\n"
	cmd += ".println \n"
	cmd += ".call firstname"

	ss.sorlOrchestration(cmd, &alp)

	//fmt.Println(alp["_func.name.firstname"])

}

func callTest12() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .tag\n"

	cmd += ".println tag name tag1\n"
	cmd += ".var sr:tags=tag\n"
	cmd += "# start the tag tag1\n"
	cmd += ".tag tag1,tagone {\n"
	cmd += "   .println Hello, I as in tag 'tag1,tagone\n"
	cmd += "   .if 10 == 10 {\n"
	cmd += "      .println 10 is 10\n"
	cmd += "      .println good..and done.\n"
	cmd += "   }\n"
	cmd += "   .println After IF inside TAG\n"
	cmd += "}\n"
	cmd += ".println Done tag..\n"
	cmd += ".println \n"

	ss.sorlOrchestration(cmd, &alp)

	//fmt.Println(alp["_func.name.firstname"])

}

func callTest13() {
	ss := &SorlSSH{}
	alp := Property{}

	alp["he.range.list"] = "Helo\n1one\n2-two-2\n33-three-333\n4-fortyfour"

	cmd := "#Testcase for .range\n"

	cmd += ".println range name \n"

	cmd += "# start the range\n"
	cmd += ".range  {he.range.list} {\n"
	cmd += "   .println Hello, I as in range={range.value}\n"
	cmd += "   .if 10 == 10 {\n"
	cmd += "      .println 10 is 10\n"
	cmd += "      .println good..and done.\n"
	cmd += "   }\n"
	cmd += "   .println After IF inside TAG\n"
	cmd += "   .println\n"
	cmd += "}\n"
	cmd += ".println Done tag..\n"
	cmd += ".println \n"

	ss.sorlOrchestration(cmd, &alp)

	//fmt.Println(alp["_func.name.firstname"])

}

func callTest14() {
	ss := &SorlSSH{}
	alp := Property{}

	alp["he.range.list"] = "Helo\n1one\n2-two-2\n33-three-333\n4-fortyfour"

	cmd := "#Testcase for .load\n"
	cmd += ".var load.file.name=/tmp/main_load.sorl\n"
	cmd += ".println .load test \n"
	cmd += "# start the range\n"
	cmd += ".load {load.file.name}\n"
	cmd += ".println Done load..\n"
	cmd += ".println \n"

	ss.sorlOrchestration(cmd, &alp)

	//fmt.Println(alp["_func.name.firstname"])

}

func callTest15() {
	ss := &SorlSSH{}
	alp := Property{}

	alp["he.range.list"] = "Helo\n1one\n2-two-2\n33-three-333\n4-fortyfour"

	cmd := "#Testcase for .load\n"
	cmd += ".var load.file.name=/tmp/main_load.sorl\n"
	cmd += ".println .load test \n"
	cmd += "# start the range\n"
	cmd += ".load {load.file.name}\n"
	cmd += ".println Done load..\n"
	cmd += ".println \n"

	ss.sorlOrchestration(cmd, &alp)

	//fmt.Println(alp["_func.name.firstname"])

}

func callTest16() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .load\n"
	cmd += ".var aa=This is a go scripting for go goers\n"
	cmd += ".var rep=go\n"
	cmd += ".var new=golang hjh;l\n"
	cmd += ".replace me.new.str {aa} {rep} {new}\n"
	cmd += ".name {aa}\n"
	cmd += ".name {me.new.str}\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest17() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .while\n"
	cmd += ".println Welcome to .while command\n"
	cmd += ".var c.i=1\n"
	cmd += ".while  10 != 10 || {c.i} != 11 {\n"
	cmd += "   .name {c.i}\n"
	cmd += "   .incr {c.i}\n"
	cmd += "   .sleep 1\n"
	cmd += "}\n"
	cmd += ".println\n"
	cmd += ".println While is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest18() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .animate\n"
	cmd += ".println Welcome to .animate command\n"
	cmd += ".animate 44 Hello World, Welcome to Google's golang and Solution Orchestration Language\n"
	cmd += ".println\n"
	cmd += ".println Animate is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func main() {

	/*
		callTest1()
		callTest2()
		callTest3()
		callTest4()
		callTest5()
		callTest6()
		callTest7()
		callTest8()
		callTest9()
		callTest10()
		callTest11()
		callTest12()
		callTest13()
		callTest14()
		callTest15()
		callTest16()
		callTest17()


	*/

	//fmt.Println(evalCondition("10 != 9 && 5 != 5 || abc == abc && a123 == a121"))
	//fmt.Println(evalCondition("10 == 10 || a == b  ||  1 == 2 || a == a && a1 == a1 && b1 == b1"))
	callTest17()

	return

	testProp := Property{
		"a":    "111",
		"b":    "222",
		"ab":   "1122",
		"test": "Soooful",
	}

	if false {
		//oLine := "abcd{hello}{world} {a{b{c{ d  }}}}this is a prop replace{name}{ lname    	}{doit{howtodo}}"
		oLine := "{test},one={a} and two={b}, chk the prop replace{ab}"
		mLine, err1 := replaceProp(oLine, testProp)
		checkError(err1)

		fmt.Printf("\n%s\n%s", oLine, mLine)
	}
	//.Println("Hello Pavana..Welcome to Go!")
	//return

	cliArgsMap := getCliArgs()
	//fmt.Println(cliArgsMap)
	printVersion()
	versionOk := strings.TrimSpace(cliArgsMap["version"])

	if versionOk == "true" {
		return
	}

	envMap := getEnvlist([]string{"USER", "HOME", "AVA"})
	//printMap("ENVIRONMENT", map[string]string(envMap))
	//logit("\n")

	homePath := envMap["HOME"]
	userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	globalOrchFilePath := strings.TrimSpace(cliArgsMap["orchfile"])
	parallelOk := strings.TrimSpace(cliArgsMap["parallel"])

	scProp := SorlConfigProperty{}
	svMap := SorlMap{}

	sorlLoadConfigFiles(&scProp, homePath, userConfigFilePath)
	scProp.printConfig()

	err := sorlLoadGlobalVars(homePath, &svMap)

	if err != nil {
		logit(fmt.Sprintf("\ninfo: %v", err))
	}
	printMap("Global Vars", svMap)

	hostList, err := sorlProcessCliArgs(scProp, cliArgsMap)

	if err != nil {
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		return
	}
	PrintList("All the selected hosts", hostList)
	//os.Exit(1)
	sorlStart(parallelOk, globalOrchFilePath, scProp, hostList, cliArgsMap, svMap)

	fmt.Println()
	fmt.Println()
}

func printVersion() {
	fmt.Println()
	fmt.Println("SORL: Solution ORchestration Language, the .(dot) scripting")
	fmt.Println("Version: 0.1-beta, build-1.0, 13JUN2019")
	time.Sleep(3 * time.Second)
	fmt.Println()
}
