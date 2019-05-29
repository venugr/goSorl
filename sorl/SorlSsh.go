package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var wg = sync.WaitGroup{}
var mut = sync.RWMutex{}

func sshPrint(color, prn string) {
	mut.Lock()
	fmt.Print(ClrUnColor + color + prn + ClrUnColor)
	mut.Unlock()
}

func sorlParallelSsh(userName, userPasswd, hostName string, portNum int, userSshKeyFile string) (*ssh.Session, *ssh.Client, error) {

	sshConfig := configSsh(hostName, userName, userPasswd, userSshKeyFile)
	client, err := dialSsh(hostName, portNum, sshConfig)

	if err != nil {
		fmt.Printf("error: Host: %s, User: %s, Port: %v, Failed to create a session: %s\n", hostName, userName, portNum, err)
		return nil, nil, err
	}

	session, err := createSSHSession(client)

	if err != nil {
		fmt.Printf("error: Host: %s, User: %s, Port: %v, Failed to create a session: %s\n", hostName, userName, portNum, err)
		return nil, nil, err
	}

	return session, client, nil

}

func runParallelSsh(userName, userPasswd, hostName string, portNum int, userSshKeyFile string) {

	sshConfig := configSsh(hostName, userName, userPasswd, userSshKeyFile)
	client, err := dialSsh(hostName, portNum, sshConfig)

	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		os.Exit(-1)
	}

	session, err := createSSHSession(client)

	if err != nil {
		fmt.Errorf("failed to create a session: %s", err)
		os.Exit(-1)
	}

	defer session.Close()
	defer client.Close()

	//runCmd(session)
	runShell(session)

	wg.Done()
}

func runShellCmd(cmd string, sshIn io.WriteCloser) {
	//fmt.Println("in runShellCmd..->" + cmd + "<-")

	//fmt.Println("***************************************")
	_, err := sshIn.Write([]byte(cmd + "\r"))
	checkError(err)
	//fmt.Println("exit runShellCmd..")

}

func waitFor(color string, display string, waitStr []string, sshOut io.Reader) (int, string) {

	//fmt.Println("in waitFor..")
	cmdOut := ""
	cmdBuf := make([]byte, 1024)
	tempOut := ""
	breakOk := false
	waitStrMatch := -1

	for {
		//fmt.Println("sshOut Read..")
		//PrintList("CMD", waitStr)

		breakOk = false
		n, err := sshOut.Read(cmdBuf)

		//fmt.Println(n)

		if err == nil {
			tempOut = string(cmdBuf[:n])
			if !strings.HasPrefix(display, "clear") {
				sshPrint(color, tempOut)
			}
			cmdOut += tempOut
		} else {
			break
		}

		for idx, wStr := range waitStr {
			//fmt.Println("Wait Str:->" + wStr + "<-")
			//fmt.Println("->" + strings.TrimSpace(tempOut) + "<-")
			if strings.HasSuffix(strings.TrimSpace(tempOut), wStr) {
				//fmt.Println("break..")
				breakOk = true
				waitStrMatch = idx
				break
			}
		}

		if breakOk {
			break
		}
	}

	//sshPrint(color, cmdOut)
	//fmt.Println(cmdOut)
	//fmt.Println("exit waitFor..")
	if strings.HasPrefix(display, "clear") {
		sshPrint(color, cmdOut)
	}

	return waitStrMatch, cmdOut

}

func setShell(session *ssh.Session) (io.Reader, io.WriteCloser, error) {

	//fmt.Println("in setShell..")
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 80, modes); err != nil {
		log.Fatal(err)
	}

	sshOut, sshOutErr := session.StdoutPipe()
	if sshOutErr != nil {
		return nil, nil, sshOutErr
	}

	sshIn, sshInErr := session.StdinPipe()
	if sshInErr != nil {
		return nil, nil, sshInErr
	}

	shellErr := session.Shell()
	if shellErr != nil {
		return nil, nil, shellErr
	}

	return sshOut, sshIn, nil

}

func runShell(session *ssh.Session) error {

	commands := []string{
		"",
		"uname -a",
		"sleep 5",
		"pwd",
		"PS1=Venu-Go-SSH-[BMUI]\\$\\ ",
		"whoami",
		"echo 'bye'",
		"ls -l /tmp",
		"sqlplus --h",
		"sleep 2",
		"env | sort ",
		"sleep 2",
		"sqlplus /nolog @/media/common/db/versions",
		"exit",
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	/*
		osFile, _ := os.Create("/tmp/dellme.dell")
		session.Stdout = osFile
		session.Stderr = osFile
	*/

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	fmt.Println()
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)

		if err != nil {
			return err
		}

	}
	fmt.Println()

	err = session.Wait()
	if err != nil {
		return err
	}

	return nil

}

func runCmdOld(session *ssh.Session) error {

	var b bytes.Buffer
	session.Stdout = &b

	fmt.Println("Inside runCmd...")
	err := session.Run("source .bash_profile ; env | sort; ls -l $BANNER_HOME; ls -l /tmp")

	if err != nil {
		return err
	}

	fmt.Println(b.String())

	return nil

}

func createSSHSession(client *ssh.Client) (*ssh.Session, error) {

	session, err := client.NewSession()
	//fmt.Println("Inside createSSH...")
	if err != nil {
		return nil, err
	}

	return session, nil
}

func dialSsh(hostName string, portNum int, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {

	//fmt.Println("Inside run dialSSH Cmd...")
	serName := hostName + ":" + strconv.Itoa(portNum)
	fmt.Print("\nConnecting to ..." + serName)
	client, err := ssh.Dial("tcp", serName, sshConfig)

	if err != nil {
		return nil, err
	}

	fmt.Print("\nConnected to ..." + serName)
	fmt.Println()
	return client, nil
}

func configSsh(hostName, userName, userPasswd string, userSshKeyFile string) *ssh.ClientConfig {

	//sshPasswd := ssh.Password(userPasswd)
	lKey, _ := configSshKey(userSshKeyFile)

	sshConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPasswd),
			ssh.PublicKeys(lKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return sshConfig
}

func configSshKey(userSshKeyFile string) (ssh.Signer, error) {

	buf, err := ioutil.ReadFile(userSshKeyFile)

	if err != nil {
		return nil, err
	}

	lKey, err := ssh.ParsePrivateKey(buf)

	if err != nil {
		return lKey, err
	}

	return lKey, err

}
