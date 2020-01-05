package gogh

import (
	"context"
	"io"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	IOContext
	GitHubContext
	LogContext
	Root() []string
	PrimaryRoot() string
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

// GitHubContext holds configurations and environments for logging.
type LogContext interface {
	LogLevel() string
	LogFlags() int // log.Lxxx flags
	LogDate() bool
	LogTime() bool
	LogMicroSeconds() bool
	LogLongFile() bool
	LogShortFile() bool
	LogUTC() bool
}
