package internal

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// commandRunner will run a command
var commandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

func execCommand(cmd *exec.Cmd) error {
	log.Println("debug: Calling", cmd.Args[0], strings.Join(cmd.Args[1:], " "))

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
