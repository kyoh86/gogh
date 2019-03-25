package gogh

import (
	"context"
	"io"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	Stdout() io.Writer
	Stderr() io.Writer
	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
	LogLevel() string
	Root() []string
	PrimaryRoot() string
}

type implContext struct {
	context.Context
	stdout      io.Writer
	stderr      io.Writer
	gitHubUser  string
	gitHubToken string
	gitHubHost  string
	logLevel    string
	root        []string
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

func (c *implContext) Root() []string {
	return c.root
}

func (c *implContext) PrimaryRoot() string {
	return c.root[0]
}
