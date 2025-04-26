package github

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Khan/genqlient/graphql"
	github "github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v3/infra/githubv4"
	"golang.org/x/oauth2"
)

type genuineAdaptor struct {
	gqlClient   graphql.Client
	restClient  *github.Client
	host        string
	tokenSource oauth2.TokenSource
}

func (a *genuineAdaptor) GetHost() string {
	return a.host
}

const (
	DefaultHost    = "github.com"
	DefaultAPIHost = "api.github.com"
)

const ClientID = "Ov23li6aEWIxek6F8P5L"

func OAuth2Config(host string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:       fmt.Sprintf("https://%s/login/oauth/authorize", host),
			TokenURL:      fmt.Sprintf("https://%s/login/oauth/access_token", host),
			DeviceAuthURL: fmt.Sprintf("https://%s/login/device/code", host),
		},
		Scopes: []string{string(github.ScopeRepo), string(github.ScopeDeleteRepo)},
	}
}

func NewAdaptor(ctx context.Context, host string, token *Token) (Adaptor, error) {
	var source oauth2.TokenSource
	if token != nil {
		source = oauth2.ReuseTokenSource(token, &tokenSource{ctx: ctx, host: host, token: token})
	}
	if host == DefaultHost || host == DefaultAPIHost {
		return newGenuineAdaptor(ctx, DefaultHost, source), nil
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
	return newGenuineEnterpriseAdaptor(
		ctx,
		host,
		baseRESTURL.String(),
		uploadRESTURL.String(),
		baseGQLURL.String(),
		source,
	)
}

type tokenSource struct {
	ctx   context.Context
	host  string
	token *oauth2.Token
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	if s.token.Valid() {
		return s.token, nil
	}
	newToken, err := refreshAccessToken(s.ctx, s.host, s.token)
	if err != nil {
		return nil, err
	}
	s.token = newToken
	return newToken, nil
}

func refreshAccessToken(ctx context.Context, host string, token *oauth2.Token) (*oauth2.Token, error) {
	oauthConfig := OAuth2Config(host)
	tokenSource := oauthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: token.RefreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func newGenuineAdaptor(ctx context.Context, host string, tokenSource oauth2.TokenSource) Adaptor {
	httpClient := oauth2.NewClient(ctx, tokenSource)
	return &genuineAdaptor{
		host:        host,
		tokenSource: tokenSource,
		restClient:  github.NewClient(httpClient),
		gqlClient:   graphql.NewClient("https://"+DefaultAPIHost+"/graphql", httpClient),
	}
}

func newGenuineEnterpriseAdaptor(
	ctx context.Context,
	host string,
	baseRESTURL, uploadRESTURL, baseGQLURL string,
	tokenSource oauth2.TokenSource,
) (Adaptor, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)
	restClient, err := github.NewClient(httpClient).WithEnterpriseURLs(baseRESTURL, uploadRESTURL)
	if err != nil {
		return nil, err
	}
	return &genuineAdaptor{
		host:        host,
		tokenSource: tokenSource,
		restClient:  restClient,
		gqlClient:   graphql.NewClient(baseGQLURL, httpClient),
	}, nil
}

func (a *genuineAdaptor) GetAccessToken() (string, error) {
	token, err := a.tokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
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

func (a *genuineAdaptor) GetAuthenticatedUser(ctx context.Context) (*User, *Response, error) {
	return a.restClient.Users.Get(ctx, "")
}

func (a *genuineAdaptor) RepositoryList(
	ctx context.Context,
	opts *RepositoryListOptions,
) ([]*RepositoryFragment, PageInfoFragment, error) {
	repos, err := githubv4.ListRepos(
		ctx,
		a.gqlClient,
		opts.Limit,
		opts.After,
		opts.IsFork,
		opts.Privacy,
		opts.OwnerAffiliations,
		opts.OrderBy,
		opts.IsArchived,
	)
	if err != nil {
		return nil, PageInfoFragment{}, err
	}
	ingrepos := make([]*RepositoryFragment, 0, len(repos.Viewer.Repositories.Edges))
	for _, edge := range repos.Viewer.Repositories.Edges {
		f := edge.Node.RepositoryFragment
		ingrepos = append(ingrepos, &f)
	}
	return ingrepos, repos.Viewer.Repositories.PageInfo.PageInfoFragment, nil
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

func (a *genuineAdaptor) OrganizationList(
	ctx context.Context,
) ([]*Organization, *Response, error) {
	return a.restClient.Organizations.List(ctx, "", &github.ListOptions{
		PerPage: 100,
	})
}

func (a *genuineAdaptor) MemberOf(
	ctx context.Context,
	org string,
) (bool, error) {
	orgs, _, err := a.OrganizationList(ctx)
	if err != nil {
		return false, err
	}
	for _, o := range orgs {
		if *o.Login == org {
			return true, nil
		}
	}
	return false, nil
}
