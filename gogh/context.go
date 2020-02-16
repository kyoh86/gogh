package gogh

import (
	"context"
)

// Context holds configurations and environments
type Context interface {
	Root() []string
	PrimaryRoot() string

	GitHubUser() string
	GitHubToken() string
	GitHubHost() string

	context.Context
}
