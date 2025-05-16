package workspace

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// ListOptions is the options for repository search
type ListOptions struct {
	// Query is the search query
	Query string
	// Limit is the maximum number of repositories to return
	// If 0, all repositories will be returned
	Limit int
}

// FinderService is a service for searching repositories
type FinderService interface {
	// FindByReference retrieves a repository by its reference
	FindByReference(ctx context.Context, ws WorkspaceService, reference repository.Reference) (*repository.Location, error)

	// FindByPath retrieves a repository by its path
	FindByPath(ctx context.Context, ws WorkspaceService, path string) (*repository.Location, error)

	// ListAllRepository retrieves a list of repositories under all roots
	ListAllRepository(context.Context, WorkspaceService, ListOptions) iter.Seq2[*repository.Location, error]

	// ListRepositoryInRoot retrieves a list of repositories under a root
	ListRepositoryInRoot(context.Context, LayoutService, ListOptions) iter.Seq2[*repository.Location, error]
}
