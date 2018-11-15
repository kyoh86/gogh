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

// A GitHubRepository represents a GitHub repository. Impliments Remote.
type GitHubRepository struct {
	url *url.URL
}

func (r *GitHubRepository) URL() *url.URL {
	return r.url
}

func (r *GitHubRepository) IsValid() bool {
	if strings.HasPrefix(r.url.Path, "/blog/") {
		return false
	}

	// must be /{user}/{project}/?
	pathComponents := strings.Split(strings.TrimRight(r.url.Path, "/"), "/")
	if len(pathComponents) != 3 {
		return false
	}

	return true
}

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
