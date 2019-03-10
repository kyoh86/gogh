package gogh

import (
	"bytes"
	"context"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	Stdout() io.Writer
	Stderr() io.Writer
	UserName() string
	GitHubToken() string
	GitHubHost() string
	LogLevel() string
	Roots() []string
	PrimaryRoot() string
	GHEHosts() []string
}

// CurrentContext get current context from OS envars and Git configurations
func CurrentContext(ctx context.Context) (Context, error) {
	userName, err := getUserName()
	if err != nil {
		return nil, err
	}
	gitHubToken, err := getGitHubToken()
	if err != nil {
		return nil, err
	}
	gitHubHost, err := getGitHubHost()
	if err != nil {
		return nil, err
	}
	logLevel, err := getLogLevel()
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
		Context:     ctx,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		userName:    userName,
		gitHubToken: gitHubToken,
		gitHubHost:  gitHubHost,
		logLevel:    logLevel,
		roots:       roots,
		gheHosts:    gheHosts,
	}, nil
}

type implContext struct {
	context.Context
	stdout      io.Writer
	stderr      io.Writer
	userName    string
	gitHubToken string
	gitHubHost  string
	logLevel    string
	roots       []string
	gheHosts    []string
}

func (c *implContext) Stdout() io.Writer {
	return c.stdout
}

func (c *implContext) Stderr() io.Writer {
	return c.stderr
}

func (c *implContext) UserName() string {
	return c.userName
}

func (c *implContext) GitHubToken() string {
	return c.gitHubToken
}

func (c *implContext) GitHubHost() string {
	return c.gitHubHost
}

func (c *implContext) LogLevel() string {
	return c.logLevel
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

func getConf(required bool, envName, confName string, altEnvNames ...string) (string, error) {
	if val := os.Getenv(envName); val != "" {
		return val, nil
	}
	val, err := getGitConf(confName)
	if err != nil {
		return "", err
	}
	if val != "" {
		return val, nil
	}

	for _, n := range altEnvNames {
		if val := os.Getenv(n); val != "" {
			return val, nil
		}
	}
	if required {
		// Make the error if it does not match any pattern
		return "", fmt.Errorf("set %s to your gitconfig", confName)
	}
	return "", nil
}

func getGitHubToken() (string, error) {
	return getConf(false, envGoghGitHubToken, "gogh.github.token", envGitHubToken)
}

func getGitHubHost() (string, error) {
	return getConf(false, envGoghGitHubHost, "gogh.github.host", envGitHubHost)
}

func getUserName() (string, error) {
	return getConf(true, envGoghGitHubUser, "gogh.github.user", envGitHubUser, envUserName)
}

func getLogLevel() (string, error) {
	if ll := os.Getenv(envLogLevel); ll != "" {
		return ll, nil
	}
	ll, err := getGitConf("gogh.logLevel")
	if err != nil {
		return "", err
	}
	if ll != "" {
		return ll, nil
	}
	return "warn", nil // default: warn
}

func getRoots() ([]string, error) {
	var roots []string
	envRoot := os.Getenv(envRoot)
	if envRoot != "" {
		roots = filepath.SplitList(envRoot)
	}
	if len(roots) == 0 {
		rts, err := getGitConfs("gogh.root")
		if err != nil {
			return nil, err
		}
		roots = rts
	}
	if len(roots) == 0 {
		roots = []string{filepath.Join(build.Default.GOPATH, "src")}
	}

	for i, v := range roots {
		path := filepath.Clean(v)
		_, err := os.Stat(path)
		switch {
		case err == nil:
			roots[i], err = filepath.EvalSymlinks(path)
			if err != nil {
				return nil, err
			}
		case os.IsNotExist(err):
			roots[i] = path
		default:
			return nil, err
		}
	}

	return unique(roots), nil
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

var gitArgsForTest []string
var gitStdinForTest []byte

// output invokes 'git config' and handles some errors properly.
func output(args ...string) (string, error) {
	param := append(append([]string{"config", "--path", "--null"}, gitArgsForTest...), args...)
	cmd := exec.Command("git", param...)
	cmd.Stderr = os.Stderr
	if gitStdinForTest != nil {
		cmd.Stdin = bytes.NewBuffer(gitStdinForTest)
	}

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
