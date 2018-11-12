package git

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// GetOneConf fetches single git-config variable.
// returns an empty string and no error if no variable is found with the given key.
func GetOneConf(key string) (string, error) {
	return output("--get", key)
}

// GetAllConf fetches git-config variable of multiple values.
func GetAllConf(key string) ([]string, error) {
	value, err := output("--get-all", key)
	if err != nil {
		return nil, err
	}

	// No results found, return an empty slice
	if value == "" {
		return nil, nil
	}

	return strings.Split(value, "\000"), nil
}

// output invokes 'git config' and handles some errors properly.
func output(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"config", "--path", "--null"}, args...)...)
	cmd.Stderr = os.Stderr

	buf, err := cmd.Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				// The key was not found, do not treat as an error
				return "", nil
			}
		}

		return "", err
	}

	return strings.TrimRight(string(buf), "\000"), nil
}
