package gogh

import (
	"context"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	userName := getUserName()
	gitHubToken := getGitHubToken()
	gitHubHost := getGitHubHost()
	logLevel := getLogLevel()
	gheHosts := getGHEHosts()
	roots, err := getRoots()
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

func getConf(envNames ...string) string {
	for _, n := range envNames {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}

func getGitHubToken() string {
	return getConf(envGoghGitHubToken, envGitHubToken)
}

func getGitHubHost() string {
	return getConf(envGoghGitHubHost, envGitHubHost)
}

func getUserName() string {
	name := getConf(envGoghGitHubUser, envGitHubUser, envUserName)
	if name == "" {
		// Make the error if it does not match any pattern
		panic(fmt.Errorf("set %s to your environment variable", envGoghGitHubUser))
	}
	return name
}

func getLogLevel() string {
	if ll := os.Getenv(envLogLevel); ll != "" {
		return ll
	}
	return "warn" // default: warn
}

func getRoots() ([]string, error) {
	var roots []string
	envRoot := os.Getenv(envRoot)
	if envRoot == "" {
		gopaths := filepath.SplitList(build.Default.GOPATH)
		roots = make([]string, 0, len(gopaths))
		for _, gopath := range gopaths {
			roots = append(roots, filepath.Join(gopath, "src"))
		}
	} else {
		roots = filepath.SplitList(envRoot)
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

func getGHEHosts() []string {
	return unique(strings.Split(os.Getenv(envGHEHosts), " "))
}
