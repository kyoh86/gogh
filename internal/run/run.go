package run

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Run command with through in/out
func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return execCommand(cmd)
}

// Silently runs command without output
func Silently(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	return execCommand(cmd)
}

// InDir runs a command within the directory
func InDir(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	return execCommand(cmd)
}

// CommandRunner will run a command
var CommandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

func execCommand(cmd *exec.Cmd) error {
	log.Println(cmd.Args[0], strings.Join(cmd.Args[1:], " "))

	err := CommandRunner(cmd)
	if err != nil {
		return &Error{cmd, err}
	}

	return nil
}

// Error holds command and its result
type Error struct {
	Command   *exec.Cmd
	ExecError error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Command.Path, e.ExecError)
}
