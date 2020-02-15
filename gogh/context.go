package gogh

import (
	"context"
	"io"
)

// Context holds configurations and environments
type Context interface {
	Root() []string
	PrimaryRoot() string

	GitHubContext
	IOContext

	context.Context
}

// IOContext holds configurations and environments for I/O.
type IOContext interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

// GitHubContext holds configurations and environments for GitHub access.
type GitHubContext interface {
	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
}
