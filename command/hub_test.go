package command_test

import (
	"net/url"

	"github.com/google/go-github/v29/github"
)

//go:generate interfacer -for github.com/kyoh86/gogh/internal/hub.Client -as command.HubClient -o hub.go
//go:generate mockgen -source hub.go -destination hub_mock_test.go -package command_test

func createNewRepoWithURL(u *url.URL) *github.Repository {
	ust := u.String()
	return &github.Repository{HTMLURL: &ust}
}
