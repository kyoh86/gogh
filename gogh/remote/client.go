package remote

import (
	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"golang.org/x/oauth2"
)

// NewClient builds GitHub Client with GitHub API token that is configured.
func NewClient(ctx gogh.Context) *github.Client {
	token := ctx.GitHubToken()
	if token == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	// TODO: support GHE
	// Use GITHUB_HOST like `hub`.
	// See: https://hub.github.com/hub.1.html#github-enterprise
	return client
}
