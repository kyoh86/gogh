package command_test

import (
	"net/url"

	"github.com/google/go-github/v32/github"
)

func createNewRepoWithURL(u *url.URL) *github.Repository {
	ust := u.String()
	return &github.Repository{HTMLURL: &ust}
}
