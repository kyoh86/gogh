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
	LogFlags() int // log.Lxxx flags
	Root() []string
	PrimaryRoot() string
}
