package github

import (
	"context"

	github "github.com/google/go-github/v33/github"
)

type Repository = github.Repository
type RepositoryCreateForkOptions = github.RepositoryCreateForkOptions
type RepositoryListByOrgOptions = github.RepositoryListByOrgOptions
type RepositoryListOptions = github.RepositoryListOptions
type Response = github.Response
type User = github.User

type Adaptor interface {
	UserGet(ctx context.Context, user string) (*User, *Response, error)

	RepositoryList(ctx context.Context, user string, opts *RepositoryListOptions) ([]*Repository, *Response, error)
	RepositoryListByOrg(ctx context.Context, org string, opts *RepositoryListByOrgOptions) ([]*Repository, *Response, error)

	RepositoryCreate(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error)
	RepositoryCreateFork(ctx context.Context, owner string, repo string, opts *RepositoryCreateForkOptions) (*Repository, *Response, error)

	RepositoryDelete(ctx context.Context, owner string, repo string) (*Response, error)
}
