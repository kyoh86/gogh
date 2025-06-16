package hook_update

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

// Execute applies a new script to a hook identified by its ID.
func (uc *UseCase) Execute(ctx context.Context, hookID, name, useCase, event, repoPattern string, content io.Reader) error {
	found, err := uc.hookService.Get(ctx, hookID)
	if err != nil {
		return err
	}
	if name != "" {
		found.Name = name
	}
	if useCase != "" {
		found.Target.UseCase = hook.UseCase(useCase)
	}
	if event != "" {
		found.Target.TriggerEvent = hook.Event(event)
	}
	if repoPattern != "" {
		found.Target.RepoPattern = repoPattern
	}
	found.UpdatedNow()
	return uc.hookService.Update(ctx, *found, content)
}
