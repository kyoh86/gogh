package context

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"time"
)

type MockContext struct {
	MStdin           io.Reader
	MStdout          io.Writer
	MStderr          io.Writer
	MGitHubUser      string
	MGitHubToken     string
	MGitHubHost      string
	MLogLevel        string
	MLogFlags        int
	MLogDate         bool
	MLogTime         bool
	MLogMicroSeconds bool
	MLogLongFile     bool
	MLogShortFile    bool
	MLogUTC          bool
	MRoot            []string
}

var _ context.Context = &MockContext{}

// Deadline : empty context.Context
func (*MockContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done : empty context.Context
func (*MockContext) Done() <-chan struct{} {
	return nil
}

// Err : empty context.Context
func (*MockContext) Err() error {
	return nil
}

// Value : empty context.Context
func (*MockContext) Value(key interface{}) interface{} {
	return nil
}

func (*MockContext) String() string {
	return "gogh/internal/context.MockContext"
}

func (c *MockContext) Stdin() io.Reader {
	if r := c.MStdin; r != nil {
		return r
	}
	return &bytes.Buffer{}
}

func (c *MockContext) Stdout() io.Writer {
	if w := c.MStdout; w != nil {
		return w
	}
	return ioutil.Discard
}

func (c *MockContext) Stderr() io.Writer {
	if w := c.MStderr; w != nil {
		return w
	}
	return ioutil.Discard
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

func (c *MockContext) LogDate() bool         { return c.MLogDate }
func (c *MockContext) LogTime() bool         { return c.MLogTime }
func (c *MockContext) LogMicroSeconds() bool { return c.MLogMicroSeconds }
func (c *MockContext) LogLongFile() bool     { return c.MLogLongFile }
func (c *MockContext) LogShortFile() bool    { return c.MLogShortFile }
func (c *MockContext) LogUTC() bool          { return c.MLogUTC }

func (c *MockContext) Root() []string {
	return c.MRoot
}

func (c *MockContext) PrimaryRoot() string {
	return c.MRoot[0]
}
