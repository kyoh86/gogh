package context

import (
	"context"
	"io"
)

type MockContext struct {
	context.Context
	MStdin        io.Reader
	MStdout       io.Writer
	MStderr       io.Writer
	MGitHubUser   string
	MGitHubToken  string
	MGitHubHost   string
	MLogLevel     string
	MLogFlags     int
	MLogDate      bool
	MLogTime      bool
	MLogLongFile  bool
	MLogShortFile bool
	MLogUTC       bool
	MRoot         []string
}

func (c *MockContext) Stdin() io.Reader {
	return c.MStdin
}

func (c *MockContext) Stdout() io.Writer {
	return c.MStdout
}

func (c *MockContext) Stderr() io.Writer {
	return c.MStderr
}

func (c *MockContext) GitHubUser() string {
	return c.MGitHubUser
}

func (c *MockContext) GitHubToken() string {
	return c.MGitHubToken
}

func (c *MockContext) GitHubHost() string {
	return c.MGitHubHost
}

func (c *MockContext) LogLevel() string {
	return c.MLogLevel
}

func (c *MockContext) LogFlags() int {
	return c.MLogFlags
}

func (c *MockContext) LogDate() bool      { return c.MLogDate }
func (c *MockContext) LogTime() bool      { return c.MLogTime }
func (c *MockContext) LogLongFile() bool  { return c.MLogLongFile }
func (c *MockContext) LogShortFile() bool { return c.MLogShortFile }
func (c *MockContext) LogUTC() bool       { return c.MLogUTC }

func (c *MockContext) Root() []string {
	return c.MRoot
}

func (c *MockContext) PrimaryRoot() string {
	return c.MRoot[0]
}
