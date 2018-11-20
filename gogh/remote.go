package gogh

import (
	"fmt"
	"net/url"
	"strings"
)

// RemoteRepo represents a remote repository.
type RemoteRepo interface {
	// The repository URL.
	URL() *url.URL
	// Checks if the URL is valid.
	IsValid() bool
}

// GitHubRepository represents a GitHub repository which implements a interface `RemoteRepo`.
type GitHubRepository struct {
	url *url.URL
}

// URL for the GitHub repository
func (r *GitHubRepository) URL() *url.URL {
	return r.url
}

// IsValid checks is the valid GitHub repository identifier valid
func (r *GitHubRepository) IsValid() bool {
	if strings.HasPrefix(r.url.Path, "/blog/") {
		return false
	}

	// must be /{user}/{project}/?
	pathComponents := strings.Split(strings.TrimRight(r.url.Path, "/"), "/")
	return len(pathComponents) == 3
}

// NewRepository create remote repository identifier from the url
func NewRepository(ctx Context, url *url.URL) (RemoteRepo, error) {
	if url.Host == "github.com" {
		return &GitHubRepository{url}, nil
	}

	gheHosts := ctx.GHEHosts()

	for _, host := range gheHosts {
		if url.Host == host {
			return &GitHubRepository{url}, nil
		}
	}

	return nil, fmt.Errorf("not supported host: %q", url.Host)
}
