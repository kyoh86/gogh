package githubv4

import (
	"context"
	"net/url"

	"github.com/machinebox/graphql"
)

type genuineAdaptor struct {
	client *graphql.Client
}

type Option func(gqlOpts *[]graphql.ClientOption, endpoint *url.URL)

func WithScheme(scheme string) Option {
	return func(gqlOpts *[]graphql.ClientOption, endpoint *url.URL) {
		endpoint.Scheme = scheme
	}
}

func WithHost(host string) Option {
	if host == DefaultHost || host == DefaultAPIHost {
		return func(*[]graphql.ClientOption, *url.URL) {}
	}
	return func(gqlOpts *[]graphql.ClientOption, endpoint *url.URL) {
		endpoint.Host = host
	}
}

func WithToken(ctx context.Context, token string) Option {
	return func(gqlOpts *[]graphql.ClientOption, endpoint *url.URL) {
		client := NewAuthClient(ctx, token)
		*gqlOpts = append(*gqlOpts, graphql.WithHTTPClient(client))
	}
}

func NewAdaptor(options ...Option) genuineAdaptor {
	var gqlOpts []graphql.ClientOption
	u := &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/graphql",
	}
	for _, opt := range options {
		opt(&gqlOpts, u)
	}
	return &genuineAdaptor{graphql.NewClient(u.String(), gqlOpts...)}
}

type genuineAdaptor interface {
	// UserGet(ctx context.Context, user string) (*User, *Response, error)

	RepositoryList(ctx context.Context, user string, opts *RepositoryListOptions) ([]*Repository, *Response, error)
	// RepositoryListByOrg(ctx context.Context, org string, opts *RepositoryListByOrgOptions) ([]*Repository, *Response, error)

	// RepositoryCreate(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error)
	// RepositoryCreateFork(ctx context.Context, owner string, repo string, opts *RepositoryCreateForkOptions) (*Repository, *Response, error)
	// RepositoryCreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *TemplateRepoRequest) (*Repository, *Response, error)
	// RepositoryDelete(ctx context.Context, owner string, repo string) (*Response, error)
	// RepositoryGet(ctx context.Context, owner string, repo string) (*Repository, *Response, error)
}
