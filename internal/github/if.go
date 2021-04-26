package github

import (
	"context"

	github "github.com/google/go-github/v35/github"
)

type (
	Repository                  = github.Repository
	RepositoryCreateForkOptions = github.RepositoryCreateForkOptions
	TemplateRepoRequest         = github.TemplateRepoRequest
	ListOptions                 = github.ListOptions
	Response                    = github.Response
	User                        = github.User
)

type Adaptor interface {
	UserGet(ctx context.Context, user string) (*User, *Response, error)
	RepositoryCreate(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error)
	RepositoryCreateFork(ctx context.Context, owner string, repo string, opts *RepositoryCreateForkOptions) (*Repository, *Response, error)
	RepositoryCreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *TemplateRepoRequest) (*Repository, *Response, error)
	RepositoryDelete(ctx context.Context, owner string, repo string) (*Response, error)
	RepositoryGet(ctx context.Context, owner string, repo string) (*Repository, *Response, error)
}
