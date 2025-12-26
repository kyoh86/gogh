package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v80/github"
	"github.com/kyoh86/gogh/v4/core/auth"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type Token = oauth2.Token

type AuthenticateService struct{}

func NewAuthenticateService() *AuthenticateService {
	return &AuthenticateService{}
}

func (s *AuthenticateService) Authenticate(ctx context.Context, host string, verify auth.Verify) (string, *Token, error) {
	config := &oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:       fmt.Sprintf("https://%s/login/oauth/authorize", host),
			TokenURL:      fmt.Sprintf("https://%s/login/oauth/access_token", host),
			DeviceAuthURL: fmt.Sprintf("https://%s/login/device/code", host),
		},
		Scopes: []string{string(github.ScopeRepo), string(github.ScopeDeleteRepo)},
	}
	// Request device code
	deviceCodeResp, err := config.DeviceAuth(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("requesting device code: %w", err)
	}
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := verify(egCtx, auth.DeviceAuthResponse{
			VerificationURI: deviceCodeResp.VerificationURI,
			UserCode:        deviceCodeResp.UserCode,
		}); err != nil {
			return fmt.Errorf("verifying device code: %w", err)
		}
		return nil
	})

	// Poll for token
	var token *Token
	eg.Go(func() error {
		// copy deviceCodeResp to avoid conflict with the other goroutine
		codeResp := *deviceCodeResp
		codeResp.Interval++ // Add a second for safety; the server may not be ready yet
		resp, err := config.DeviceAccessToken(egCtx, &codeResp)
		if err != nil {
			return fmt.Errorf("polling for token: %w", err)
		}
		token = resp
		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", nil, err
	}

	if token == nil {
		return "", nil, fmt.Errorf("got nil token response")
	}

	// Get user info
	conn := getConnection(ctx, host, token)
	user, _, err := conn.rest.Users.Get(ctx, "")
	if err != nil {
		return "", nil, fmt.Errorf("getting authenticated user: %w", err)
	}
	return user.GetLogin(), token, nil
}

var _ auth.AuthenticateService = (*AuthenticateService)(nil)
