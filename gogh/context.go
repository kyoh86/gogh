package gogh

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	LogLevel() string
	UserName() string
	Roots() []string
	PrimaryRoot() string
	GHEHosts() []string
}

// CurrentContext get current context from OS envars and Git configurations
func CurrentContext(ctx context.Context) (Context, error) {
	logLevel, err := getLogLevel()
	if err != nil {
		return nil, err
	}
	userName, err := getUserName()
	if err != nil {
		return nil, err
	}
	roots, err := getRoots()
	if err != nil {
		return nil, err
	}
	gheHosts, err := getGHEHosts()
	if err != nil {
		return nil, err
	}
	return &implContext{
		Context:  ctx,
		logLevel: logLevel,
		userName: userName,
		roots:    roots,
		gheHosts: gheHosts,
	}, nil
}

type implContext struct {
	context.Context
	logLevel string
	userName string
	roots    []string
	gheHosts []string
}

func (c *implContext) LogLevel() string {
	return c.logLevel
}

func (c *implContext) UserName() string {
	return c.userName
}

func (c *implContext) Roots() []string {
	return c.roots
}

func (c *implContext) PrimaryRoot() string {
	rts := c.Roots()
	return rts[0]
}

func (c *implContext) GHEHosts() []string {
	return c.gheHosts
}

func getLogLevel() (string, error) {
	user, err := getGitConf("gogh.log")
	if err != nil {
		return "", err
	}
	if user != "" {
		return user, nil
	}
	if user := os.Getenv("GOGH_LOG_LEVEL"); user != "" {
		return user, nil
	}
	return "Info", nil
}

func getUserName() (string, error) {
	user, err := getGitConf("gogh.user")
	if err != nil {
		return "", err
	}
	if user != "" {
		return user, nil
	}
	if user := os.Getenv("GITHUB_USER"); user != "" {
		return user, nil
	}
	switch runtime.GOOS {
	case "windows":
		if user := os.Getenv("USERNAME"); user != "" {
			return user, nil
		}
	default:
		if user := os.Getenv("USER"); user != "" {
			return user, nil
		}
	}
	// Make the error if it does not match any pattern
	return "", fmt.Errorf("set gogh.user to your gitconfig")
}

func getRoots() ([]string, error) {
	envRoot := os.Getenv("GOGH_ROOT")
	if envRoot != "" {
		return filepath.SplitList(envRoot), nil
	}
	rts, err := getGitConfs("gogh.root")
	if err != nil {
		return nil, err
	}

	if len(rts) == 0 {
		rts = []string{filepath.Join(build.Default.GOPATH, "src")}
	}

	for i, v := range rts {
		path := filepath.Clean(v)
		_, err := os.Stat(path)
		switch {
		case err == nil:
			rts[i], err = filepath.EvalSymlinks(path)
			if err != nil {
				return nil, err
			}
		case os.IsNotExist(err):
			rts[i] = path
		default:
			return nil, err
		}
	}

	return rts, nil
}

// PrimaryRoot returns the first one of the root directories to clone repository.
func getPrimaryRoot() (string, error) {
	rts, err := getRoots()
	if err != nil {
		return "", err
	}
	return rts[0], nil
}

func getGHEHosts() ([]string, error) {
	return getGitConfs("gogh.ghe.host")
}

// getGitConf fetches single git-config variable.
// returns an empty string and no error if no variable is found with the given key.
func getGitConf(key string) (string, error) {
	return output("--get", key)
}

// getGitConfs fetches git-config variable of multiple values.
func getGitConfs(key string) ([]string, error) {
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
	cmd := exec.Command("git", append([]string{"config", "--type", "path", "--null"}, args...)...)
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
