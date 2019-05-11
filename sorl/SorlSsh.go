package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

var wg = sync.WaitGroup{}

func sorlParallerlSsh(userName, userPasswd, hostName string, portNum int) (*ssh.Session, *ssh.Client, error) {

	sshConfig := configSsh(hostName, userName, userPasswd)
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

func runParallelSsh(userName, userPasswd, hostName string, portNum int) {

	sshConfig := configSsh(hostName, userName, userPasswd)
	client, err := dialSsh(hostName, portNum, sshConfig)

	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		os.Exit(-1)
	}

	session, err := createSSHSession(client)

	if err != nil {
		fmt.Errorf("Failed to create a session: %s\n", err)
		os.Exit(-1)
	}

	defer session.Close()
	defer client.Close()

	//runCmd(session)
	runShell(session)

	wg.Done()
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

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	fmt.Println("\n")
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)

		if err != nil {
			return err
		}

	}
	fmt.Println("\n")

	err = session.Wait()
	if err != nil {
		return err
	}

	return nil

}

func runCmd(session *ssh.Session) error {

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
	fmt.Println("Inside createSSH...")
	if err != nil {
		return nil, err
	}

	return session, nil
}

func dialSsh(hostName string, portNum int, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {

	fmt.Println("Inside run dialSSH Cmd...")
	serName := hostName + ":" + strconv.Itoa(portNum)
	fmt.Println(serName)
	client, err := ssh.Dial("tcp", serName, sshConfig)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func configSsh(hostName, userName, userPasswd string) *ssh.ClientConfig {

	sshConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(userPasswd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return sshConfig
}
