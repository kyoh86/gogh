package github

import (
	"context"
	"net/http"

	github "github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type genuineAdaptor struct {
	client *github.Client
}

func (c *genuineAdaptor) UserGet(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return c.client.Users.Get(ctx, user)
}
func (c *genuineAdaptor) RepositoryList(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return c.client.Repositories.List(ctx, user, opts)
}
func (c *genuineAdaptor) RepositoryListByOrg(ctx context.Context, org string, opts *github.RepositoryListByOrgOptions) ([]*github.Repository, *github.Response, error) {
	return c.client.Repositories.ListByOrg(ctx, org, opts)
}
func (c *genuineAdaptor) RepositoryCreate(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return c.client.Repositories.Create(ctx, org, repo)
}
func (c *genuineAdaptor) RepositoryCreateFork(ctx context.Context, owner string, repo string, opts *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error) {
	return c.client.Repositories.CreateFork(ctx, owner, repo, opts)
}
func (c *genuineAdaptor) RepositoryDelete(ctx context.Context, owner string, repo string) (*github.Response, error) {
	return c.client.Repositories.Delete(ctx, owner, repo)
}

func NewAuthClient(ctx context.Context, accessToken string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
}

func NewAdaptor(httpClient *http.Client) Adaptor {
	return &genuineAdaptor{
		client: github.NewClient(httpClient),
	}
}

func NewEnterpriseAdaptor(ctx context.Context, baseURL string, uploadURL string, httpClient *http.Client) (Adaptor, error) {
	client, err := github.NewEnterpriseClient(baseURL, uploadURL, httpClient)
	if err != nil {
		return nil, err
	}
	return &genuineAdaptor{
		client: client,
	}, nil
}
