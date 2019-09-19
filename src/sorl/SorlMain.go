package main

import (
	"fmt"
	"strings"
	"time"
)

var sorlDebug = false

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
	cmd += "   .if 5 <= 10 {\n"
	cmd += "      .println 5 <= 10 \n"
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
	cmd += "df -h\n"
	cmd += ".if 5 > 5 && 1 <= 10 {\n"
	cmd += "   .println Multi if conditon works!\n"
	cmd += "   .if 5 >= 8 {\n"
	cmd += "      .println Inner If works\n"
	cmd += "   }\n"
	cmd += "   .println Hello IF done.\n"
	cmd += "}\n"
	cmd += ".name DONE DUNA DONE\n"

	ss.sorlOrchestration(cmd, &alp)

	fmt.Println(alp["c.z"])
}

func callTest5_1() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for If\n"
	cmd += ".println Test case for if\n"
	cmd += "\n"
	cmd += ".var exp.v=1 != 2 && 1 == 1 && 5 == 5 && 4 <= 1\n"
	cmd += ".println exp={exp.v}\n"
	cmd += ".if {exp.v} {\n"
	cmd += "   .name It Work's\n"
	cmd += "}\n"

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

func callTest19() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .while\n"
	cmd += ".println Welcome to .while command\n"
	cmd += ".var c.i=1\n"
	cmd += ".while  10 != 10 || {c.i} != 11 {\n"
	cmd += "   .var lw.i=0\n"
	cmd += "   .name {c.i}\n"
	cmd += "   .while {lw.i} != {c.i} {\n"
	cmd += "      .println lw.i={lw.i}\n"
	cmd += "      .incr {lw.i}\n"
	cmd += "   }\n"
	cmd += "   .incr {c.i}\n"
	cmd += "#   .sleep 1\n"
	cmd += "}\n"
	cmd += ".println\n"
	cmd += ".println While is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest20() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .while\n"
	cmd += ".println Welcome to .while command\n"
	cmd += ".var c.i=1\n"
	cmd += ".while  {c.i} < 10 || {c.i} <= 10 {\n"
	cmd += "   .name {c.i}\n"
	cmd += "   .incr {c.i}\n"
	cmd += "   .sleep 1\n"
	cmd += "}\n"
	cmd += ".println\n"
	cmd += ".println While is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest21() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .style colors\n"
	cmd += ".println Welcome to .style command\n"
	cmd += ".var st.txt=Welcome to SORL, the .(dot) Scripting\n"
	cmd += ".println\n"
	cmd += ".style red {st.txt}\n"
	cmd += ".style white {st.txt}\n"
	cmd += ".style green {st.txt}\n"
	cmd += ".println\n"
	cmd += ".style blue {st.txt}\n"
	cmd += ".style yellow {st.txt}\n"
	cmd += ".style magenta {st.txt}\n"
	cmd += ".style cyan {st.txt}\n"

	cmd += ".println Styling is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest22() {

	ss := &SorlSSH{}
	alp := Property{}
	alp["test.str"] = "1 2 3 4 5 6 7 8 9 10\n"
	alp["test.str"] += "one two three four five six seven eight nine ten\n"
	alp["test.str"] += "One    Two Three Four Five Six      Seven Eight Nine Ten\n"
	alp["test.str"] += "One1 Two Three3 Four Five5 Six Seven7 Eight Nine9 Ten\n"
	alp["test.str"] += "aa bb ccc 444 Five 666\n"
	alp["test.str"] += "ONE TWO THREE FOUR FIVE      SIX SEVEN EIGHT  NINE TEN\n"

	cmd := "#Testcase for .select \n"
	cmd += ".println Welcome to .select command\n"
	cmd += ".select 9 from {test.str}\n"
	cmd += ".println {_select.result.str}\n"
	cmd += ".println Select is done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest23() {

	ss := &SorlSSH{}
	alp := Property{}
	alp["test.str"] = "1 2 3 4 5 6 7 8 9 10\n"
	alp["test.str"] += "one two three four five six seven eight nine ten\n"
	alp["test.str"] += "One    Two Three Four Five Six      Seven Eight Nine Ten\n"
	alp["test.str"] += "One1 Two Three3 Four Five5 Six Seven7 Eight Nine9 Ten\n"
	alp["test.str"] += "aa bb ccc 444 Five 666\n"
	alp["test.str"] += "ONE TWO THREE FOUR FIVE      SIX SEVEN EIGHT  NINE TEN"

	cmd := "#Testcase for .trimleft .trimright\n"
	cmd += ".println \n"
	cmd += ".println Welcome to .trimleft .trimright command\n"
	cmd += ".println ==>{test.str}<==\n"
	cmd += ".println\n"
	cmd += ".trimleft  new.str {test.str}\n"
	cmd += ".println\n"
	cmd += ".println ==>{new.str}<==\n"
	cmd += ".trimleft  new.str {new.str}\n"
	cmd += ".println\n"
	cmd += ".println =*>{new.str}<*=\n"
	cmd += ".trimright  new.str {new.str}\n"
	cmd += ".println\n"
	cmd += ".println =RIG>{new.str}<RIG=\n"
	cmd += ".println \n"
	cmd += ".println Trimleft Trimright is done.\n"
	cmd += ".println \n"

	ss.sorlOrchestration(cmd, &alp)
}

func callTest24() {

	repeat := "Hello, Welome to Google Golang and Solution Orchestration Language"

	for {
		for _, ch := range repeat {
			fmt.Print(string(ch))
			time.Sleep(time.Millisecond * 50)
		}
		fmt.Print("\r")
		time.Sleep(time.Millisecond * 100)
		for range repeat {
			fmt.Print(" ")
		}
		time.Sleep(time.Millisecond * 20)
		fmt.Print("\r")
	}
}

func callTest25() {
	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .incr by/.decr by \n"
	cmd += ".println name .incr by/.decr by\n"
	cmd += ".input sl.no1 Enter Emp Number? \n"
	cmd += ".println\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".incr {sl.no1} by 2\n.incr {sl.no1} by 3\n"
	cmd += ".println Num={sl.no1}\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".incr {sl.no1}\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".decr {sl.no1}\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".decr {sl.no1}     by    5\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".decr {sl.no1} by 2\n"
	cmd += ".name {sl.no1}\n"
	cmd += ".println Done with incr by/decr by\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest26() {
	ss := &SorlSSH{}
	alp := Property{}

	alp["_cmd.output"] = "This is also erased..."

	cmd := "#Testcase for .clear \n"
	cmd += ".println name .clear\n"
	cmd += ".var cl.txt=Hello This will be cleared.\n"
	cmd += ".name {cl.txt}\n"
	cmd += ".clear {cl.txt}\n"
	cmd += ".name {cl.txt}\n"
	cmd += ".name {_cmd.output}\n"
	cmd += ".clear\n"
	cmd += ".name {_cmd.output}\n"
	cmd += ".println Done with .clear\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest27() {
	ss := &SorlSSH{}
	alp := Property{}

	alp["_cmd.output"] = "This is also erased..."

	cmd := "#Testcase for .set \n"
	cmd += ".println Test for .set\n"
	cmd += ".var cl.txt=Hello This will be reset.\n"
	cmd += ".name {cl.txt}\n"
	cmd += ".set cl.txt=How is this, after being reset.\n"
	cmd += ".name {cl.txt}\n"
	cmd += ".println Done with .set\n"

	ss.sorlOrchestration(cmd, &alp)

}

func callTest28() {

	ss := &SorlSSH{}
	alp := Property{}

	cmd := "#Testcase for .read a file\n"
	cmd += ".println\n"
	cmd += ".var hello.file.name=/Users/venul/Development/go/code/src/github.com/goRecodeSorl/goSorl/LICENSE\n"
	cmd += ".println read file {hello.file.name}\n"
	cmd += ".read hello.txt {hello.file.name}\n"
	cmd += ".range {hello.txt} {\n"
	cmd += "   .println \n"
	cmd += "   .print {range.value}\n"
	cmd += "}\n"
	cmd += ".write {hello.txt} {hello.file.name}.write\n"
	cmd += ".println Done.\n"

	ss.sorlOrchestration(cmd, &alp)
}

func newMain() {

	scProp := SorlConfigProperty{}
	svMap := SorlMap{}

	cliArgsMap := getCliArgs()

	actName, actValue, err := sorlGetAction(cliArgsMap)
	if err != nil {
		return
	}

	if actName == "encrypt" || actName == "decrypt" {
		sorlEncDecCliArgs(scProp, cliArgsMap)
		return
	}

	printVersion()

	versionOk := strings.TrimSpace(cliArgsMap["version"])
	if versionOk == "true" {
		return
	}

	envMap := getEnvlist([]string{"USER", "HOME", "SORL_DEBUG", "AVA"})

	if strings.ToLower(envMap["SORL_DEBUG"]) == "yes" {
		sorlDebug = true
	}

	homePath := envMap["HOME"]
	userConfigFilePath := strings.TrimSpace(cliArgsMap["config"])
	//globalOrchFilePath := strings.TrimSpace(cliArgsMap["orchfile"])
	//parallelOk := strings.TrimSpace(cliArgsMap["parallel"])
	varFileName := strings.TrimSpace(cliArgsMap["var-file"])

	sorlLoadConfigFiles(&scProp, homePath, userConfigFilePath)

	actArgs, err := sorlGetActionArgs(actName, scProp, cliArgsMap)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	//fmt.Println(actArgs)

	err = sorlLoadGlobalVars(homePath, &svMap)
	if err != nil {
		logit(fmt.Sprintf("\ninfo: %v", err))
	}

	if varFileName != "" {
		err = sorlLoadFileVars(varFileName, &svMap)

		if err != nil {
			logit(fmt.Sprintf("\ninfo: %v", err))
		}
	}

	err = sorlArgsVars(&svMap)

	if err != nil {
		logit(fmt.Sprintf("\ninfo: %v", err))
	}

	if sorlDebug {
		printMap("Global/Loaded Vars", svMap)
	}

	if actName == "conn" {
		sorlActionConn(actName, actValue, actArgs, cliArgsMap, svMap, scProp)
	}

}

func main() {

	newMain()
	fmt.Println()
	return

	/*
		key := "123456789012345678901234"
		by := sorlEncryptText(key, "u_pick_it")
		fmt.Println(string(by))

		by = sorlDecryptText(key, by)
		fmt.Println(string(by))

		return
	*/

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
	//callTest28()

	/*
		color.Cyan.Printf("Simple to use %s\n", "color")
		color.Error.Println("message")
		color.Info.Prompt("message")
		color.Warn.Println("message")
	*/

	//return

	/*

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

	*/

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

	ok := sorlEncDecCliArgs(scProp, cliArgsMap)

	if ok {
		return
	}

	/*
		ok = sorlConnectCliArgs(scProp, cliArgsMap)

		if ok {
			return
		}
	*/

	sorlLoadConfigFiles(&scProp, homePath, userConfigFilePath)
	scProp.printConfig()

	err := sorlLoadGlobalVars(homePath, &svMap)

	if err != nil {
		logit(fmt.Sprintf("\ninfo: %v", err))
	}

	err = sorlArgsVars(&svMap)

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
	currentTime := time.Now()
	fmt.Println()
	fmt.Println("SORL: Solution ORchestration Language, the .(dot) scripting")
	fmt.Println("Version: 0.1-beta, build-1.0-17-SEP-2019, " + currentTime.Format("02-Jan-06"))
	time.Sleep(1 * time.Second)
	fmt.Println()
}
