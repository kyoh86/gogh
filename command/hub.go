package command

import (
	"net/url"

	"github.com/kyoh86/gogh/command/internal"
	"github.com/kyoh86/gogh/gogh"
)

type hubClient interface {
	Fork(ctx gogh.Context, project *gogh.Project, repo *gogh.Repo, noRemote bool, remoteName, organization string) error
	Create(ctx gogh.Context, project *gogh.Project, repo *gogh.Repo, description string, homepage *url.URL, private, browse, clipboard bool) error
}

type mockHubClient struct {
}

func (i *mockHubClient) Fork(gogh.Context, *gogh.Project, *gogh.Repo, bool, string, string) error {
	return nil
}

func (i *mockHubClient) Create(gogh.Context, *gogh.Project, *gogh.Repo, string, *url.URL, bool, bool, bool) error {
	return nil
}

var defaultHubClient hubClient = &internal.HubClient{}

func hub() hubClient {
	return defaultHubClient
}
