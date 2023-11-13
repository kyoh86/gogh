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
	host       string
	restClient *github.Client
	gqlClient  *githubv4.Client
}

func (a *genuineAdaptor) GetHost() string {
	return a.host
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
		return newGenuineAdaptor(DefaultHost, client), nil
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
	return newGenuineEnterpriseAdaptor(
		host,
		baseRESTURL.String(),
		uploadRESTURL.String(),
		baseGQLURL.String(),
		client,
	)
}

func NewAuthClient(ctx context.Context, accessToken string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
}

func newGenuineAdaptor(host string, httpClient *http.Client) Adaptor {
	return &genuineAdaptor{
		host:       host,
		restClient: github.NewClient(httpClient),
		gqlClient:  githubv4.NewClient(httpClient, "https://"+DefaultAPIHost+"/graphql"),
	}
}

func newGenuineEnterpriseAdaptor(
	host string,
	baseRESTURL, uploadRESTURL, baseGQLURL string,
	httpClient *http.Client,
) (Adaptor, error) {
	restClient, err := github.NewEnterpriseClient(baseRESTURL, uploadRESTURL, httpClient)
	if err != nil {
		return nil, err
	}
	return &genuineAdaptor{
		host:       host,
		restClient: restClient,
		gqlClient:  githubv4.NewClient(httpClient, baseGQLURL),
	}, nil
}

func (a *genuineAdaptor) UserGet(ctx context.Context, user string) (*User, *Response, error) {
	return a.restClient.Users.Get(ctx, user)
}

func (a *genuineAdaptor) GetMe(
	ctx context.Context,
) (string, error) {
	me, _, err := a.restClient.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}
	return me.GetLogin(), nil
}

func (a *genuineAdaptor) RepositoryList(
	ctx context.Context,
	opts *RepositoryListOptions,
) ([]*RepositoryFragment, PageInfoFragment, error) {
	repos, err := a.gqlClient.ListRepos(
		ctx,
		opts.Limit,
		opts.After,
		opts.IsFork,
		opts.Privacy,
		opts.OwnerAffiliations,
		opts.OrderBy,
	)
	if err != nil {
		return nil, PageInfoFragment{}, err
	}
	ingrepos := make([]*RepositoryFragment, 0, len(repos.Viewer.Repositories.Edges))
	for _, edge := range repos.Viewer.Repositories.Edges {
		ingrepos = append(ingrepos, edge.Node)
	}
	return ingrepos, repos.Viewer.Repositories.PageInfo, nil
}

func (a *genuineAdaptor) RepositoryCreate(
	ctx context.Context,
	org string,
	repo *Repository,
) (*Repository, *Response, error) {
	return a.restClient.Repositories.Create(ctx, org, repo)
}

func (a *genuineAdaptor) RepositoryCreateFork(
	ctx context.Context,
	owner string,
	repo string,
	opts *RepositoryCreateForkOptions,
) (*Repository, *Response, error) {
	return a.restClient.Repositories.CreateFork(ctx, owner, repo, opts)
}

func (a *genuineAdaptor) RepositoryGet(
	ctx context.Context,
	owner string,
	repo string,
) (*Repository, *Response, error) {
	return a.restClient.Repositories.Get(ctx, owner, repo)
}

func (a *genuineAdaptor) RepositoryDelete(
	ctx context.Context,
	owner string,
	repo string,
) (*Response, error) {
	return a.restClient.Repositories.Delete(ctx, owner, repo)
}

func (a *genuineAdaptor) RepositoryCreateFromTemplate(
	ctx context.Context,
	templateOwner, templateRepo string,
	templateRepoReq *TemplateRepoRequest,
) (*Repository, *Response, error) {
	return a.restClient.Repositories.CreateFromTemplate(
		ctx,
		templateOwner,
		templateRepo,
		templateRepoReq,
	)
}
