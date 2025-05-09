package workspace

import (
	"context"
	"iter"
)

// ListOptions is the options for repository search
type ListOptions struct {
	// Query is the search query
	Query string
	// Limit is the maximum number of repositories to return
	// If 0, all repositories will be returned
	Limit int
}

// Repository is informations of a repository
type RepoInfo interface {
	// FullPath returns the full path of the repository
	FullPath() string
	// Path returns the path from root of the repository
	Path() string
	// Host returns the host of the repository
	Host() string
	// Owner returns the owner of the repository
	Owner() string
	// Name returns the name of the repository
	Name() string
}

// FinderService is a service for searching repositories
type FinderService interface {
	// FindByPath retrieves a repository by its reference
	FindByPath(ctx context.Context, ws WorkspaceService, path string) (RepoInfo, error)

	// ListAllRepository retrieves a list of repositories under all roots
	ListAllRepository(context.Context, WorkspaceService, ListOptions) iter.Seq2[RepoInfo, error]

	// ListRepositoryInRoot retrieves a list of repositories under a root
	ListRepositoryInRoot(context.Context, LayoutService, ListOptions) iter.Seq2[RepoInfo, error]
}
