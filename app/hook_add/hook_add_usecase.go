package hook_add

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/hook"
)

type Options struct {
	Name        string
	RepoPattern string
	UseCase     string
	Event       string
}

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

func (uc *UseCase) Execute(ctx context.Context, opts Options, content io.Reader) (*hook.Hook, error) {
	h := hook.Hook{
		ID:   uuid.NewString(),
		Name: opts.Name,
		Target: hook.Target{
			RepoPattern: opts.RepoPattern,
			UseCase:     hook.UseCase(opts.UseCase),
			Event:       hook.Event(opts.Event),
		},
	}
	h.CreatedNow()
	return &h, uc.hookService.Add(ctx, h, content)
}
