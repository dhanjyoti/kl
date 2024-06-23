package sshclient

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHConfig struct {
	Host    string
	User    string
	KeyPath string
	SSHPort int
}

var ErrSSHNotReady = fmt.Errorf("ssh is not ready")

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return ssh.PublicKeys(key), nil
}

func DoSSH(sc SSHConfig) error {
	pkFile, err := publicKeyFile(sc.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %s, please ensure you have the correct key", err)
	}

	config := &ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			pkFile,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sc.Host, sc.SSHPort), config)
	if err != nil {
		return ErrSSHNotReady
	}
	defer client.Close()

	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s, please try again", err)
	}
	defer session.Close()

	// Create a session

	// Allocate a pseudo-terminal (pty) for the session
	ptmx, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to create pseudo-terminal: %s, please try again", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), ptmx)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	go handleResize(session)

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height = 100, 30
	}

	// Start the session with a pseudo-terminal
	if err := session.RequestPty("xterm", height, width, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("failed to start pseudo-terminal: %s, please try again", err)
	}

	// Start the remote shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %s, please try again", err)

	}

	// Wait for the session to finish
	if err := session.Wait(); err != nil {
		term.Restore(int(os.Stdin.Fd()), ptmx)
		return nil
		//fn.Warnf("session exited with error: %s", err.Error())
	}

	return nil
}
