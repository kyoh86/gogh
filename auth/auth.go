// Package auth defines GitHub authorization
package auth

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/util"
	"github.com/octokit/go-octokit/octokit"
)

var (
	requiredScopes = map[string]bool{"public_repo": true, "repo": true, "user": true}
)

type scopeError struct {
	required map[string]bool
}

func (s *scopeError) Error() string {
	return strings.Join(util.StringBoolMapKeys(s.required), ",")
}

func checkScopes(scopes []string) error {
	required := map[string]bool{}
	for scope := range requiredScopes {
		required[scope] = true
	}
	for _, scope := range scopes {
		delete(required, scope)
	}
	if len(required) == 0 {
		return nil
	}
	return &scopeError{required: required}
}

// NewClient will connect to GitHub with the access-token, and return client
func NewClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	_, res, err := client.Users.Get("me")
	if err != nil {
		return nil, err
	}

	if err := checkScopes(strings.Split(res.Header.Get("X-OAuth-Scopes"), ", ")); err != nil {
		return nil, fmt.Errorf("failed to get client: GOGH needs scopes [%s]", err.Error())
	}
	return client, nil
}

// Login will connect to GitHub with username and password, and return an access-token
func Login(user, password string) (token string, err error) {
	client := octokit.NewClient(octokit.BasicAuth{Login: user, Password: password})
	u, _ := octokit.AuthorizationsURL.Expand(nil)
	authorization, res := client.Authorizations(u).Create(octokit.M{
		"scopes":      util.StringBoolMapKeys(requiredScopes),
		"note":        fmt.Sprintf("%s; %s", env.AppName, env.AppDescription),
		"note_url":    env.SiteURL,
		"fingerprint": time.Now().Format(time.RFC3339Nano),
	})
	if res.HasError() {
		return "", res
	}
	return authorization.Token, nil
}
