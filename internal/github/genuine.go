package github

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-github/v35/github"
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

type Option func(baseRESTURL *url.URL, uploadRESTURL *url.URL, baseGQLURL *url.URL)

func WithScheme(scheme string) Option {
	return func(baseRESTURL *url.URL, uploadRESTURL *url.URL, baseGQLURL *url.URL) {
		baseRESTURL.Scheme = scheme
		uploadRESTURL.Scheme = scheme
		baseGQLURL.Scheme = scheme
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
	baseRESTURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/v3",
	}
	uploadRESTURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/uploads",
	}
	baseGQLURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/graphql",
	}
	for _, option := range options {
		option(baseRESTURL, uploadRESTURL, baseGQLURL)
	}
	return newGenuineEnterpriseAdaptor(baseRESTURL.String(), uploadRESTURL.String(), baseGQLURL.String(), client)
}

func NewAuthClient(ctx context.Context, accessToken string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
}

func newGenuineAdaptor(httpClient *http.Client) Adaptor {
	return &genuineAdaptor{
		restClient: github.NewClient(httpClient),
		gqlClient:  githubv4.NewClient(httpClient, "https://"+DefaultAPIHost+"/graphql"),
	}
}

func newGenuineEnterpriseAdaptor(baseRESTURL, uploadRESTURL, baseGQLURL string, httpClient *http.Client) (Adaptor, error) {
	restClient, err := github.NewEnterpriseClient(baseRESTURL, uploadRESTURL, httpClient)
	if err != nil {
		return nil, err
	}
	return &genuineAdaptor{
		restClient: restClient,
		gqlClient:  githubv4.NewClient(httpClient, baseGQLURL),
	}, nil
}

func (c *genuineAdaptor) UserGet(ctx context.Context, user string) (*User, *Response, error) {
	return c.restClient.Users.Get(ctx, user)
}

type RepositoryListOptions struct { // TODO:
	Limit             *int64
	Cursor            *string
	IsFork            *bool
	Privacy           *githubv4.RepositoryPrivacy
	OwnerAffiliations []*githubv4.RepositoryAffiliation
	OrderBy           *githubv4.RepositoryOrder
}

func (c *genuineAdaptor) RepositoryList(ctx context.Context, opts *RepositoryListOptions) ([]*Repository, error) {
	var ()
	repos, err := c.gqlClient.ListRepos(
		ctx,
		opts.Limit,
		opts.Cursor,
		opts.IsFork,
		opts.Privacy,
		opts.OwnerAffiliations,
		opts.OrderBy,
	)
	if err != nil {
		return nil, err
	}
	ingrepos := make([]*Repository, 0, len(repos.Viewer.Repositories.Edges))
	for _, edge := range repos.Viewer.Repositories.Edges {
		srcrepo := edge.Node
		ingrepo := &Repository{
			URL: &srcrepo.URL,
			Owner: &github.User{
				Login: &srcrepo.Owner.Login,
			},
			Name:        &srcrepo.Name,
			Description: srcrepo.Description,
			Fork:        &srcrepo.IsFork,
			Archived:    &srcrepo.IsArchived,
			Private:     &srcrepo.IsPrivate,
			IsTemplate:  &srcrepo.IsTemplate,
			// TODO: CreatedAt    string  "json:\"createdAt\" graphql:\"createdAt\""
			// TODO: PushedAt     *string "json:\"pushedAt\" graphql:\"pushedAt\""
			// TODO: Parent       *struct {
			// TODO: 	Owner struct {
			// TODO: 		ID    string "json:\"id\" graphql:\"id\""
			// TODO: 		Login string "json:\"login\" graphql:\"login\""
			// TODO: 	} "json:\"owner\" graphql:\"owner\""
			// TODO: 	Name         string "json:\"name\" graphql:\"name\""
			// TODO: } "json:\"parent\" graphql:\"parent\""
		}
		ingrepos = append(ingrepos, ingrepo)
	}
	return ingrepos, nil
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
