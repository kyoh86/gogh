package repo

import (
	"fmt"
	"net/url"
	"regexp"
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

// A GitHubGistRepository represents a GitHub Gist repository.
type GitHubGistRepository struct {
	url *url.URL
}

func (r *GitHubGistRepository) URL() *url.URL {
	return r.url
}

func (r *GitHubGistRepository) IsValid() bool {
	return true
}

type GoogleCodeRepository struct {
	url *url.URL
}

func (r *GoogleCodeRepository) URL() *url.URL {
	return r.url
}

var validGoogleCodePathPattern = regexp.MustCompile(`^/p/[^/]+/?$`)

func (r *GoogleCodeRepository) IsValid() bool {
	return validGoogleCodePathPattern.MatchString(r.url.Path)
}

type BluemixRepository struct {
	url *url.URL
}

func (r *BluemixRepository) URL() *url.URL {
	return r.url
}

var validBluemixPathPattern = regexp.MustCompile(`^/git/[^/]+/[^/]+$`)

func (r *BluemixRepository) IsValid() bool {
	return validBluemixPathPattern.MatchString(r.url.Path)
}

type OtherRepository struct {
	url *url.URL
}

func (r *OtherRepository) URL() *url.URL {
	return r.url
}

func (r *OtherRepository) IsValid() bool {
	return true
}

func NewRepository(url *url.URL) (Remote, error) {
	if url.Host == "github.com" {
		return &GitHubRepository{url}, nil
	}

	if url.Host == "gist.github.com" {
		return &GitHubGistRepository{url}, nil
	}

	if url.Host == "code.google.com" {
		return &GoogleCodeRepository{url}, nil
	}

	if url.Host == "hub.jazz.net" {
		return &BluemixRepository{url}, nil
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

	return &OtherRepository{url}, nil
}
