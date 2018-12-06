package gogh

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// run command with through in/out
func run(ctx Context, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()

	return execCommand(cmd)
}

// runSilently runs command without output
func runSilently(ctx Context, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard

	return execCommand(cmd)
}

// runInDir runs a command within the directory
func runInDir(ctx Context, dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = dir

	return execCommand(cmd)
}

// commandRunner will run a command
var commandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

func execCommand(cmd *exec.Cmd) error {
	log.Println("debug: calling", cmd.Args[0], strings.Join(cmd.Args[1:], " "))

	err := commandRunner(cmd)
	if err != nil {
		return &execError{cmd, err}
	}

	return nil
}

// execError holds command and its result
type execError struct {
	Command   *exec.Cmd
	ExecError error
}

func (e *execError) Error() string {
	return fmt.Sprintf("%s: %s", e.Command.Path, e.ExecError)
}
