package delegate

import (
	"fmt"
	"os/exec"
)

func ExecCommand(cmd *exec.Cmd) error {
	err := cmd.Run()
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
