package logout

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/auth"
)

type UseCase struct {
	tokenService auth.TokenService
}

func NewUseCase(tokenService auth.TokenService) *UseCase {
	return &UseCase{
		tokenService: tokenService,
	}
}

func (uc *UseCase) Execute(_ context.Context, host, owner string) error {
	return uc.tokenService.Delete(host, owner)
}
