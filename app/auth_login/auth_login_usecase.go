package auth_login

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
)

type UseCase struct {
	tokenService   auth.TokenService
	authService    auth.AuthenticateService
	hostingService hosting.HostingService
}

func NewUseCase(
	tokenService auth.TokenService,
	authService auth.AuthenticateService,
	hostingService hosting.HostingService,
) *UseCase {
	return &UseCase{
		tokenService:   tokenService,
		authService:    authService,
		hostingService: hostingService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, host string, verify auth.Verify) error {
	user, token, err := uc.authService.Authenticate(ctx, host, verify)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err := uc.tokenService.Set(host, user, *token); err != nil {
		return fmt.Errorf("failed to set token: %w", err)
	}
	return nil
}
