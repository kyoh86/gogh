package github

import (
	"context"

	github "github.com/google/go-github/v33/github"
)

type Adaptor interface {
	GetUser(ctx context.Context, user string) (*github.User, *github.Response, error)

	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListRepositoriesByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error)

	CreateRepository(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	CreateForkRepository(ctx context.Context, owner string, repo string, opts *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)

	DeleteRepositories(ctx context.Context, owner string, repo string) (*github.Response, error)
}
