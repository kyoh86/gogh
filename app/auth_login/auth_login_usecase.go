package auth_login

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"golang.org/x/oauth2" //TODO: remove oauth2 dependency
)

type UseCase struct {
	tokenService   auth.TokenService
	hostingService hosting.HostingService
}

func NewUseCase(tokenService auth.TokenService) *UseCase {
	return &UseCase{
		tokenService: tokenService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, host string) (chan *oauth2.DeviceAuthResponse, chan error) {
	resch := make(chan *oauth2.DeviceAuthResponse)
	errch := make(chan error)

	go uc.executeCore(ctx, resch, errch, host)

	return resch, errch
}

func (uc *UseCase) executeCore(ctx context.Context, resch chan *oauth2.DeviceAuthResponse, errch chan error, host string) {
	defer close(resch)
	defer close(errch)

	// Get OAuth2 config
	config, err := uc.hostingService.GetOauth2Config(ctx, host)
	if err != nil {
		errch <- fmt.Errorf("failed to get OAuth2 config: %w", err)
		return
	}

	// Request device code
	deviceCodeResp, err := config.DeviceAuth(ctx)
	if err != nil {
		errch <- fmt.Errorf("failed to request device code: %w", err)
		return
	}
	resch <- deviceCodeResp

	// Poll for token
	deviceCodeResp.Interval++ // Add a second for safety
	tokenResp, err := config.DeviceAccessToken(ctx, deviceCodeResp)
	if err != nil {
		errch <- fmt.Errorf("failed to poll for token: %w", err)
		return
	}

	if tokenResp == nil {
		errch <- fmt.Errorf("got nil token response")
		return
	}

	// Get user info
	user, err := uc.hostingService.GetAuthenticatedUserName(ctx, host, tokenResp)
	if err != nil {
		errch <- fmt.Errorf("failed to get authenticated user info: %w", err)
		return
	}

	if err := uc.tokenService.Set(host, user, *tokenResp); err != nil {
		errch <- fmt.Errorf("failed to save token: %w", err)
		return
	}
}
