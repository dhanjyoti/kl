package functions

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecCmd(cmdString string, env map[string]string, verbose bool) error {
	r := csv.NewReader(strings.NewReader(cmdString))
	r.Comma = ' '
	cmdArr, err := r.Read()
	if err != nil {
		return err
	}
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	if verbose {
		Log("[#] " + strings.Join(cmdArr, " "))
		cmd.Stdout = os.Stdout
	}

	// cmd.Env = os.Environ()

	if env == nil {
		env = map[string]string{}
	}

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stderr = os.Stderr
	// s.Start()
	err = cmd.Run()
	// s.Stop()
	return err
}

func Exec(cmdString string, env map[string]string) ([]byte, error) {
	r := csv.NewReader(strings.NewReader(cmdString))
	r.Comma = ' '
	cmdArr, err := r.Read()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)

	// fmt.Println(cmd.String(), cmdString)
	cmd.Stderr = os.Stderr

	if env == nil {
		env = map[string]string{}
	}

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	return out, err
}

// isAdmin checks if the current user has administrative privileges.
func isAdmin() bool {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func WinSudoExec(cmdString string, env map[string]string) ([]byte, error) {
	if isAdmin() {
		return Exec(cmdString, env)
	}

	r := csv.NewReader(strings.NewReader(cmdString))
	r.Comma = ' '
	cmdArr, err := r.Read()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)

	quotedArgs := strings.Join(cmdArr[1:], ",")

	return Exec(fmt.Sprintf("powershell -Command Start-Process -WindowStyle Hidden -FilePath %s -ArgumentList %q -Verb RunAs", cmd.Path, quotedArgs), map[string]string{"PATH": os.Getenv("PATH")})
}

func Confirm(yes string, defaultValue string) bool {
	var response string
	_, _ = fmt.Scanln(&response)
	if response == "" {
		if defaultValue == "" {
			return false
		}
		response = defaultValue
	}
	return strings.ToLower(response) == strings.ToLower(yes)
}
