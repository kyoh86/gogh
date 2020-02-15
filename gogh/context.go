package gogh

import (
	"context"
)

// Context holds configurations and environments
type Context interface {
	Root() []string
	PrimaryRoot() string

	GitHubContext

	context.Context
}

// GitHubContext holds configurations and environments for GitHub access.
type GitHubContext interface {
	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
}
