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

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

var wg = sync.WaitGroup{}
var mut = sync.RWMutex{}

func (ss SorlSSH) sshPrint(prn string, allProp *Property) {

	mut.Lock()

	mStr := ClrUnColor + ss.sorlSshColor + prn + ClrUnColor

	if sorlWindows {
		mStr = prn
	}
	fmt.Print(mStr)

	mut.Unlock()

	fileNames := (*allProp)["_log.file.names"]
	fileNames = strings.TrimSuffix(fileNames, ",")

	for _, lFileName := range strings.Split(fileNames, ",") {
		(*allProp)["_log.data."+lFileName] += prn
	}

	logVarNames := (*allProp)["_log.var.names"]
	logVarNames = strings.TrimSuffix(logVarNames, ",")

	for _, lVarName := range strings.Split(logVarNames, ",") {
		(*allProp)[lVarName] += prn
	}

}

func sshPrint(color, prn string, allProp *Property) {

	mut.Lock()

	mStr := ClrUnColor + color + prn + ClrUnColor

	if sorlWindows {
		mStr = prn
	}

	//(*allProp)["sr:echo"] = "on"
	if (*allProp)["sr:echo"] == "on" {
		fmt.Print(mStr)
	}

	mut.Unlock()

	fileNames := (*allProp)["_log.file.names"]
	fileNames = strings.TrimSuffix(fileNames, ",")

	for _, lFileName := range strings.Split(fileNames, ",") {
		(*allProp)["_log.data."+lFileName] += prn
	}

	logVarNames := (*allProp)["_log.var.names"]
	logVarNames = strings.TrimSuffix(logVarNames, ",")

	for _, lVarName := range strings.Split(logVarNames, ",") {
		(*allProp)[lVarName] += prn
	}

}

func (ss *SorlSSH) sorlParallelSsh() error {

	ss.configSsh()
	err := ss.dialSsh()

	if err != nil {
		fmt.Printf("error: Host: %s, User: %s, Port: %v, Failed to create a session: %s\n", ss.sorlSshHostName, ss.sorlSshUserName, ss.sorlSshHostPortNum, err)
		return err
	}

	err = ss.createSSHSession()

	if err != nil {
		fmt.Printf("error: Host: %s, User: %s, Port: %v, Failed to create a session: %s\n", ss.sorlSshHostName, ss.sorlSshUserName, ss.sorlSshHostPortNum, err)
		return err
	}

	return nil

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

func (ss *SorlSSH) runParallelSsh() {

	ss.configSsh()
	err := ss.dialSsh()

	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		os.Exit(-1)
	}

	err = ss.createSSHSession()

	if err != nil {
		fmt.Errorf("failed to create a session: %s", err)
		os.Exit(-1)
	}

	defer ss.sorlSshSession.Close()
	defer ss.sorlSshClient.Close()

	//runCmd(session)
	runShell(ss.sorlSshSession)

	wg.Done()
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

func (ss *SorlSSH) runShellCmd(cmd string) {

	_, err := ss.sorlSshIn.Write([]byte(cmd + "\n"))
	checkError(err)

}

func runShellCmd(cmd string, sshIn io.WriteCloser) {
	//fmt.Println("in runShellCmd..->" + cmd + "<-")

	//fmt.Println("***************************************")
	_, err := sshIn.Write([]byte(cmd + "\r"))
	checkError(err)
	//fmt.Println("exit runShellCmd..")

}

func waitFor(echoOn bool, color string, display string, waitStr []string, sshOut io.Reader) (int, string) {

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
			if echoOn && !strings.HasPrefix(display, "clear") {
				sshPrint(color, tempOut, nil)
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
	fmt.Println("echo: ", echoOn)
	if echoOn && strings.HasPrefix(display, "clear") {
		sshPrint(color, cmdOut, nil)
	}

	return waitStrMatch, cmdOut

}

func (ss *SorlSSH) waitFor(echoOn bool, display string, waitStr []string, allProp *Property) (int, string) {

	//fmt.Println("in waitFor..")

	color := ss.sorlSshColor
	sshOut := ss.sorlSshOut

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
			if echoOn && !strings.HasPrefix(display, "clear") {
				sshPrint(color, tempOut, allProp)
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
	if echoOn && strings.HasPrefix(display, "clear") {
		sshPrint(color, cmdOut, allProp)
	}

	return waitStrMatch, cmdOut

}

func (ss *SorlSSH) setShell() error {

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := ss.sorlSshSession.RequestPty("xterm", 80, 180, modes); err != nil {
		log.Fatal(err)
	}

	sshOut, sshOutErr := ss.sorlSshSession.StdoutPipe()
	if sshOutErr != nil {
		return sshOutErr
	}

	sshIn, sshInErr := ss.sorlSshSession.StdinPipe()
	if sshInErr != nil {
		return sshInErr
	}

	shellErr := ss.sorlSshSession.Shell()
	if shellErr != nil {
		return shellErr
	}

	ss.sorlSshIn = sshIn
	ss.sorlSshOut = sshOut

	return nil

}

func setShell(session *ssh.Session) (io.Reader, io.WriteCloser, error) {

	//fmt.Println("in setShell..")
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 180, modes); err != nil {
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

func (ss *SorlSSH) createSSHSession() error {

	lClient := ss.sorlSshClient

	session, err := lClient.NewSession()

	if err != nil {
		return err
	}

	ss.sorlSshSession = session

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

func (ss *SorlSSH) dialSsh() error {

	lHostIP := ss.sorlSshHostIP
	lPortNum := ss.sorlSshHostPortNum
	lUserName := ss.sorlSshUserName

	//fmt.Println("Inside run dialSSH Cmd...")
	serName := lHostIP + ":" + strconv.Itoa(lPortNum)
	fmt.Print("\nConnecting to ..." + lUserName + "@" + serName + "\n")
	client, err := ssh.Dial("tcp", serName, ss.sorlSshClientConfig)

	if err != nil {
		return err
	}

	fmt.Print("Connected to ..." + lUserName + "@" + serName)
	fmt.Println()
	ss.sorlSshClient = client
	return nil
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

func (ss *SorlSSH) configSsh() {

	lUserName := ss.sorlSshUserName
	lUserPassword := ss.sorlSshUserPassword

	lKey, err := ss.configSshKey()

	if err != nil {
		//fmt.Println("\nError: Config SSH Key issue...")
		//return

		lKey = nil
	}

	sshConfig := &ssh.ClientConfig{
		User: lUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(lUserPassword),
			ssh.PublicKeys(lKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	ss.sorlSshClientConfig = sshConfig

}

func configSsh(hostName, userName, userPasswd string, userSshKeyFile string) *ssh.ClientConfig {

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

func (ss *SorlSSH) configSshKey() (ssh.Signer, error) {

	buf, err := ioutil.ReadFile(ss.sorlSshHostKeyFile)

	if err != nil {
		return nil, err
	}

	lKey, err := ssh.ParsePrivateKey(buf)

	if err != nil {
		return lKey, err
	}

	return lKey, err
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

func (ss *SorlSSH) sorlSftp(isTemplate bool, filesSrcMap SorlMap, filesDestMap SorlMap, allProp *Property) {

	sftp, err := sftp.NewClient(ss.sorlSshClient)

	if err != nil {

	}

	defer sftp.Close()

	/*
		w := sftp.Walk("/tmp")

		for w.Step() {
			fmt.Println(w.Path())
		}

	*/

	ss.sshPrint("\n", allProp)

	for lKey, _ := range filesSrcMap {

		sftp.MkdirAll(filesDestMap[lKey])

		remFile := filesDestMap[lKey] + "/" + lKey
		locFile := filesSrcMap[lKey] + "/" + lKey

		ss.sshPrint("copying file '"+locFile+"'...\n", allProp)

		newFile, err := sftp.Create(remFile)

		if err != nil {

		}

		fullFile, err := ReadBinaryFile(locFile)
		if err != nil {

		}

		if isTemplate {
			fullString, _ := replaceProp(string(fullFile), (*allProp))
			fullFile = []byte(fullString)
		}

		if _, err := newFile.Write(fullFile); err != nil {

		}

		ss.sshPrint("copied file '"+locFile+"'...\n\n", allProp)
	}

	/*
		fStat, err := sftp.Lstat("test.sorl")
		if err != nil {

		}

		fmt.Println(fStat)
	*/
}
