package add

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

type Usecase struct {
	hookService hook.HookService
}

func NewUsecase(hookService hook.HookService) *Usecase {
	return &Usecase{hookService: hookService}
}

func (uc *Usecase) Execute(ctx context.Context, opts Options) (string, error) {
	h := hook.Entry{
		Name:          opts.Name,
		RepoPattern:   opts.RepoPattern,
		TriggerEvent:  hook.Event(opts.TriggerEvent),
		OperationType: hook.OperationType(opts.OperationType),
		OperationID:   opts.OperationID,
	}
	return uc.hookService.Add(ctx, h)
}
