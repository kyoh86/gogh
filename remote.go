package gogh

import (
	"context"
	"net/http"

	"github.com/kyoh86/gogh/v2/internal/github"
)

type RemoteController struct {
	adaptor github.Adaptor
}

func GithubAdaptor(ctx context.Context, server Server) (github.Adaptor, error) {
	var client *http.Client
	if server.Token() != "" {
		client = github.NewAuthClient(ctx, server.Token())
	}
	//UNDONE: support Enterprise with server.baseURL and server.uploadURL
	return github.NewAdaptor(client), nil
}

func NewRemoteController(adaptor github.Adaptor) *RemoteController {
	return &RemoteController{
		adaptor: adaptor,
	}
}

type RemoteListOption struct {
	User    string
	Query   string
	Options *github.RepositoryListOptions
}

func (c *RemoteController) List(ctx context.Context, option *RemoteListOption) ([]Project, error) {
	// UNDONE: implement
	return nil, nil
}

type RemoteListByOrgOption struct {
	Query   string
	Options *github.RepositoryListByOrgOptions
}

func (c *RemoteController) ListByOrg(ctx context.Context, org string, option *RemoteListByOrgOption) ([]Project, error) {
	// UNDONE: implement
	return nil, nil
}
