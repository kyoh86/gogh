package delegate

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// CommandRunner will run a command
var CommandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

func ExecCommand(cmd *exec.Cmd) error {
	log.Println("debug: Calling", cmd.Args[0], strings.Join(cmd.Args[1:], " "))

	err := CommandRunner(cmd)
	if err != nil {
		return &ExecError{cmd, err}
	}

	return nil
}

// ExecError holds command and its result
type ExecError struct {
	Command   *exec.Cmd
	ExecError error
}

func (e *ExecError) Error() string {
	return fmt.Sprintf("%s: %s", e.Command.Path, e.ExecError)
}
