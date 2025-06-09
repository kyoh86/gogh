package hook_list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

func (uc *UseCase) Execute(ctx context.Context) iter.Seq2[*hook.Hook, error] {
	return uc.hookService.ListHooks()
}
