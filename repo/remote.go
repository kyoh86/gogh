package repo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kyoh86/gogh/internal/git"
)

// A Remote represents a remote repository.
type Remote interface {
	// The repository URL.
	URL() *url.URL
	// Checks if the URL is valid.
	IsValid() bool
}

// GitHubRepository represents a GitHub repository which implements a interface `Remote`.
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
func NewRepository(url *url.URL) (Remote, error) {
	if url.Host == "github.com" {
		return &GitHubRepository{url}, nil
	}

	gheHosts, err := git.GetAllConf("gogh.ghe.host")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve GH:E hostname from .gitconfig: %s", err)
	}

	for _, host := range gheHosts {
		if url.Host == host {
			return &GitHubRepository{url}, nil
		}
	}

	return nil, fmt.Errorf("not supported host: %q", url.Host)
}
