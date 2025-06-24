package update

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type Options struct {
	Name          string
	RepoPattern   string
	TriggerEvent  string
	OperationType string
	OperationID   string
}

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

func (uc *UseCase) Execute(ctx context.Context, id string, opts Options) error {
	h := hook.Entry{
		Name:          opts.Name,
		RepoPattern:   opts.RepoPattern,
		TriggerEvent:  hook.Event(opts.TriggerEvent),
		OperationType: hook.OperationType(opts.OperationType),
		OperationID:   opts.OperationID,
	}
	return uc.hookService.Update(ctx, id, h)
}
