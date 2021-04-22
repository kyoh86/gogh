package github

import (
	"context"
	"net/http"
	"net/url"

	github "github.com/google/go-github/v35/github"
	"github.com/kyoh86/gogh/v2/internal/githubv4"
	"golang.org/x/oauth2"
)

type genuineAdaptor struct {
	restClient *github.Client
	gqlClient  *githubv4.Client
}

const (
	DefaultHost    = "github.com"
	DefaultAPIHost = "api.github.com"
)

type Option func(baseURL *url.URL, uploadURL *url.URL)

func WithScheme(scheme string) Option {
	return func(baseURL *url.URL, uploadURL *url.URL) {
		baseURL.Scheme = scheme
		uploadURL.Scheme = scheme
	}
}

func WithBasePath(path string) Option {
	return func(baseURL *url.URL, _ *url.URL) {
		baseURL.Path = path
	}
}

func WithUploadPath(path string) Option {
	return func(_ *url.URL, uploadURL *url.URL) {
		uploadURL.Path = path
	}
}

func NewAdaptor(ctx context.Context, host, token string, options ...Option) (Adaptor, error) {
	var client *http.Client
	if token != "" {
		client = NewAuthClient(ctx, token)
	}
	if host == DefaultHost || host == DefaultAPIHost {
		return newGenuineAdaptor(client), nil
	}
	baseURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/v3",
	}
	uploadURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/uploads",
	}
	for _, option := range options {
		option(baseURL, uploadURL)
	}
	return newGenuineEnterpriseAdaptor(baseURL.String(), uploadURL.String(), client)
}

func NewAuthClient(ctx context.Context, accessToken string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
}

func newGenuineAdaptor(httpClient *http.Client) Adaptor {
	return &genuineAdaptor{
		restClient: github.NewClient(httpClient),
	}
}

func newGenuineEnterpriseAdaptor(baseURL string, uploadURL string, httpClient *http.Client) (Adaptor, error) {
	client, err := github.NewEnterpriseClient(baseURL, uploadURL, httpClient)
	if err != nil {
		return nil, err
	}
	return &genuineAdaptor{
		restClient: client,
	}, nil
}

func (c *genuineAdaptor) UserGet(ctx context.Context, user string) (*User, *Response, error) {
	return c.restClient.Users.Get(ctx, user)
}

func (c *genuineAdaptor) SearchRepository(ctx context.Context, query string, opts *SearchOptions) ([]*github.Repository, *Response, error) {
	result, resp, err := c.restClient.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, resp, err
	}
	return result.Repositories, resp, nil
}

func (c *genuineAdaptor) RepositoryCreate(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error) {
	return c.restClient.Repositories.Create(ctx, org, repo)
}

func (c *genuineAdaptor) RepositoryCreateFork(ctx context.Context, owner string, repo string, opts *RepositoryCreateForkOptions) (*Repository, *Response, error) {
	return c.restClient.Repositories.CreateFork(ctx, owner, repo, opts)
}

func (c *genuineAdaptor) RepositoryGet(ctx context.Context, owner string, repo string) (*Repository, *Response, error) {
	return c.restClient.Repositories.Get(ctx, owner, repo)
}

func (c *genuineAdaptor) RepositoryDelete(ctx context.Context, owner string, repo string) (*Response, error) {
	return c.restClient.Repositories.Delete(ctx, owner, repo)
}

func (c *genuineAdaptor) RepositoryCreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *TemplateRepoRequest) (*Repository, *Response, error) {
	return c.restClient.Repositories.CreateFromTemplate(ctx, templateOwner, templateRepo, templateRepoReq)
}

func (c *genuineAdaptor) OrganizationsList(ctx context.Context, opts *github.ListOptions) ([]*Organization, *Response, error) {
	return c.restClient.Organizations.List(ctx, "", opts)
}
