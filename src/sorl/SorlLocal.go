package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func sorlLocal(parallelOk, orchFile string, scProp SorlConfigProperty,
	hostsList []string, cliArgsMap map[string]string,
	svMap map[string]string) {

	allProp := Property{}

	for fKey, fVal := range svMap {
		allProp[fKey] = fVal
	}

	cmdStr := ""

	for _, lVal := range strings.Split(allProp["_cmd.arg.order"], ",") {

		cmdStr += allProp[lVal] + "\n"

	}

	cmdStr = strings.TrimSpace(cmdStr)

	//fmt.Println("WaitPrompt: " + waitPrompt)

	commands := strings.Split(cmdStr, "\n")

	if orchFile != "" {

		commandsNew, fileErr := ReadFile(orchFile)

		if fileErr != nil {
			fmt.Print("\n\nError: ")
			fmt.Println(fileErr)
			//ss.sorlSshSession.Close()
			return
		}

		commands = commandsNew

	}

	//PrintList("CMDs", commands)

	ss := &SorlSSH{}
	allProp["_host.local"] = "yes"
	ss.sorlOrchestration(strings.Join(commands, "\n"), &allProp)

}

func SorlExecWait() {

}

func SorlExec(cmdStr string) {

	fmt.Println("\n\n=========================")
	fmt.Println("+" + cmdStr)
	fmt.Println("=========================")

	cmdOut := ""
	cmdBuf := make([]byte, 1024)
	tempOut := ""

	var wg = sync.WaitGroup{}

	//cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
	//cmd := exec.Command("go", "run", "./TestRun.go")

	//cmdName := strings.Split(cmdStr, " ")
	//cmdStr = strings.TrimLeft(cmdStr, cmdName)

	cmd := exec.Command("sh", "-c", cmdStr)

	if sorlWindows {
		cmd = exec.Command("cmd", "/C", cmdStr)
	}

	stdOut, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	defer stdOut.Close()

	stdIn, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	defer stdIn.Close()

	cmd.Stderr = os.Stderr

	wg.Add(1)
	go func() {
		if err = cmd.Start(); err != nil {
			fmt.Println("An error occured: ", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {

			n, err := stdOut.Read(cmdBuf)

			if err == nil {
				tempOut = string(cmdBuf[:n])
				cmdOut += tempOut
				fmt.Print(tempOut)
				//time.Sleep(time.Second)

				if strings.HasSuffix(strings.TrimSpace(cmdOut), "favorite number?") {
					fmt.Print("5\n")
					io.WriteString(stdIn, "5\n")
				}

				if strings.HasSuffix(strings.TrimSpace(cmdOut), "favorite number:") {
					fmt.Print("15\n")
					io.WriteString(stdIn, "15\n")
				}

				if strings.HasSuffix(strings.TrimSpace(cmdOut), "First Name=") {
					fmt.Print("VENU\n")
					io.WriteString(stdIn, "VENU\n")
				}

			} else {
				break
			}

		}
		wg.Done()
	}()

	wg.Wait()
	cmd.Wait()
}
