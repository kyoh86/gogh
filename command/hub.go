// Code generated by interfacer; DO NOT EDIT

package command

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/kyoh86/gogh/gogh"
	"net/url"
)

// HubClient is an interface generated for "github.com/kyoh86/gogh/internal/hub.Client".
type HubClient interface {
	Create(context.Context, gogh.Env, *gogh.Repo, string, *url.URL, bool) (*github.Repository, error)
	Fork(context.Context, gogh.Env, *gogh.Repo, string) (*gogh.Repo, error)
	Repos(context.Context, gogh.Env, string, bool, bool, bool, bool, string, string, string) ([]string, error)
}
