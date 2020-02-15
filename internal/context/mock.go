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

func (c *MockContext) Root() []string {
	return c.MRoot
}

func (c *MockContext) PrimaryRoot() string {
	return c.MRoot[0]
}
