package context

import (
	"context"
	"time"
)

type MockContext struct {
	MGitHubUser  string
	MGitHubToken string
	MGitHubHost  string
	MRoot        []string
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
