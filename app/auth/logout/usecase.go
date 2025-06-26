package logout

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/auth"
)

type Usecase struct {
	tokenService auth.TokenService
}

func NewUsecase(tokenService auth.TokenService) *Usecase {
	return &Usecase{
		tokenService: tokenService,
	}
}

func (uc *Usecase) Execute(_ context.Context, host, owner string) error {
	return uc.tokenService.Delete(host, owner)
}
