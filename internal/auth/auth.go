// Package auth defines GitHub authorization utilities
package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/wacul/ptr"
	"golang.org/x/oauth2"
)

var (
	requiredScopes = []github.Scope{"repo", "user"}
)

type scopeError struct {
	required []string
}

func (s *scopeError) Error() string {
	return strings.Join(s.required, ",")
}

func checkScopes(scopes []string) error {
	required := map[string]struct{}{}
	for _, scope := range requiredScopes {
		required[string(scope)] = struct{}{}
	}
	for _, scope := range scopes {
		delete(required, scope)
	}
	if len(required) == 0 {
		return nil
	}
	err := &scopeError{required: make([]string, 0, len(required))}
	for req := range required {
		err.required = append(err.required, req)
	}
	return err
}

// NewClient will connect to GitHub with the access-token, and return client
func NewClient(oAuthContext context.Context, token string) *github.Client {
	return github.NewClient(
		oauth2.NewClient(
			oAuthContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		),
	)
}

// Validate client authorization.
func Validate(ctx context.Context, client *github.Client) error {
	_, res, err := client.Users.Get(ctx, "me")
	if err != nil {
		return err
	}

	if err := checkScopes(strings.Split(res.Header.Get("X-OAuth-Scopes"), ", ")); err != nil {
		return fmt.Errorf("failed to get client: GOGH needs scopes [%s]", err.Error())
	}
	return nil
}

type basicAuthTransport struct {
	Username string
	Password string
}

func (b basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			b.Username, b.Password)))))
	return http.DefaultTransport.RoundTrip(req)
}

// Login will connect to GitHub with username and password, and return an access-token
func Login(ctx context.Context, user, password string) (token string, err error) {
	client := github.NewClient(&http.Client{Transport: &basicAuthTransport{}})
	auth, _, err := client.Authorizations.Create(ctx, &github.AuthorizationRequest{
		Scopes:      []github.Scope{"repo", "user"},
		Note:        ptr.String(fmt.Sprintf("%s; %s", "gogh", "GitHub CLI Client")),
		NoteURL:     ptr.String("https://github.com/kyoh86/gogh"),
		Fingerprint: ptr.String(time.Now().Format(time.RFC3339Nano)),
	})
	if err != nil {
		return "", err
	}
	return *auth.Token, nil
}
