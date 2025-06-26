package remove

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type Usecase struct {
	hookService hook.HookService
}

func NewUsecase(hookService hook.HookService) *Usecase {
	return &Usecase{hookService: hookService}
}

func (uc *Usecase) Execute(ctx context.Context, hookID string) error {
	return uc.hookService.Remove(ctx, hookID)
}
