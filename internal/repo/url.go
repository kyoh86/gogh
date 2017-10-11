package repo

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// GitHubURL gets url for GitHub from remotes
func (r *Repository) GitHubURL() (*url.URL, error) {
	conf, err := r.repository.Config()
	if err != nil {
		return nil, errors.Wrap(err, "get local config")
	}

	remote, ok := conf.Remotes["origin"]
	if !ok {
		return nil, errors.New("there are no origin in remotes")
	}
	return findGitHubURL(remote.URLs)
}

func findGitHubURL(urls []string) (*url.URL, error) {
	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil {
			continue
		}
		if parsedURL.Scheme == "https" && parsedURL.Host == "github.com" {
			return parsedURL, nil
		}
	}

	return nil, errors.New("there are no remote in github")
}

// Identifier for a repository in the GitHub
type Identifier struct {
	Owner  string
	Name   string
	Host   string
	Scheme string
}

// Identifier will get an identifier for a repository in the GitHub
func (r *Repository) Identifier() (*Identifier, error) {
	parsedURL, err := r.GitHubURL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get url")
	}
	return parseIdentifier(parsedURL)
}

func parseIdentifier(u *url.URL) (*Identifier, error) {
	floors := strings.Split(u.Path, "/")
	if len(floors) != 3 {
		return nil, errors.New("failed to parse url: paths not formed as owner/name")
	}
	return &Identifier{
		Owner:  floors[1],
		Name:   strings.TrimSuffix(floors[2], ".git"),
		Host:   u.Host,
		Scheme: u.Scheme,
	}, nil
}
