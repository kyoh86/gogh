package remote

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"golang.org/x/oauth2"
)

// NewClient builds GitHub Client with GitHub API token that is configured.
func NewClient(ctx gogh.Context) (*github.Client, error) {
	if host := ctx.GitHubHost(); host != "" && host != gogh.DefaultHost {
		url := fmt.Sprintf("https://%s/api/v3", host)
		return github.NewEnterpriseClient(url, url, oauth2Client(ctx, ctx.GitHubToken()))
	}

	return github.NewClient(oauth2Client(ctx, ctx.GitHubToken())), nil
}

func oauth2Client(ctx context.Context, token string) *http.Client {
	if token == "" {
		return nil
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}
