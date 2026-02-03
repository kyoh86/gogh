package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type Usecase struct {
	hookService hook.HookService
}

func NewUsecase(
	hookService hook.HookService,
) *Usecase {
	return &Usecase{
		hookService: hookService,
	}
}

func (uc *Usecase) Execute(ctx context.Context) iter.Seq2[hook.Hook, error] {
	_ = ctx
	return uc.hookService.List()
}
