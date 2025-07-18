package login

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/hosting"
)

// Usecase is the services to handle login authentication.
type Usecase struct {
	tokenService   auth.TokenService
	authService    auth.AuthenticateService
	hostingService hosting.HostingService
}

// NewUsecase creates a new Usecase instance with the provided services.
func NewUsecase(
	tokenService auth.TokenService,
	authService auth.AuthenticateService,
	hostingService hosting.HostingService,
) *Usecase {
	return &Usecase{
		tokenService:   tokenService,
		authService:    authService,
		hostingService: hostingService,
	}
}

// DeviceAuthResponse represents the response from a device authentication request.
type DeviceAuthResponse = auth.DeviceAuthResponse

// Verify is a function type to verify the authentication response.
type Verify = auth.Verify

// Execute performs the authentication process.
func (uc *Usecase) Execute(ctx context.Context, host string, verify Verify) error {
	user, token, err := uc.authService.Authenticate(ctx, host, verify)
	if err != nil {
		return fmt.Errorf("authenticating: %w", err)
	}
	if token == nil {
		return fmt.Errorf("token is nil")
	}
	if err := uc.tokenService.Set(host, user, *token); err != nil {
		return fmt.Errorf("setting token: %w", err)
	}
	return nil
}
