package gogh

import (
	"context"
	"io"
)

type implContext struct {
	context.Context
	stdout       io.Writer
	stderr       io.Writer
	gitHubUser   string
	gitHubToken  string
	gitHubHost   string
	logLevel     string
	logFlags     int
	logDate      bool
	logTime      bool
	logLongFile  bool
	logShortFile bool
	logUTC       bool
	root         []string
}

func (c *implContext) Stdout() io.Writer {
	return c.stdout
}

func (c *implContext) Stderr() io.Writer {
	return c.stderr
}

func (c *implContext) GitHubUser() string {
	return c.gitHubUser
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

func (c *implContext) LogFlags() int {
	return c.logFlags
}

func (c *implContext) LogDate() bool      { return c.logDate }
func (c *implContext) LogTime() bool      { return c.logTime }
func (c *implContext) LogLongFile() bool  { return c.logLongFile }
func (c *implContext) LogShortFile() bool { return c.logShortFile }
func (c *implContext) LogUTC() bool       { return c.logUTC }

func (c *implContext) Root() []string {
	return c.root
}

func (c *implContext) PrimaryRoot() string {
	return c.root[0]
}
