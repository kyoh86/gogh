package github

import (
	"context"

	github "github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type genuineAdaptor struct {
	client *github.Client
}

func (c *genuineAdaptor) GetUser(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return c.client.Users.Get(ctx, user)
}
func (c *genuineAdaptor) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return c.client.Repositories.List(ctx, user, opts)
}
func (c *genuineAdaptor) ListRepositoriesByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return c.client.Repositories.ListByOrg(ctx, org, opts)
}
func (c *genuineAdaptor) CreateRepository(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return c.client.Repositories.Create(ctx, org, repo)
}
func (c *genuineAdaptor) CreateForkRepository(ctx context.Context, owner string, repo string, opts *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error) {
	return c.client.Repositories.CreateFork(ctx, owner, repo, opts)
}
func (c *genuineAdaptor) DeleteRepositories(ctx context.Context, owner string, repo string) (*github.Response, error) {
	return c.client.Repositories.Delete(ctx, owner, repo)
}

func NewClient() Adaptor {
	return &genuineAdaptor{
		client: github.NewClient(nil),
	}
}

func NewAuthClient(ctx context.Context, accessToken string) Adaptor {
	return &genuineAdaptor{
		client: github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))),
	}
}
