package hook_create

import (
	"context"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/google/uuid"
	"io"
)

type Options struct {
	Name        string
	Description string
	Event       string
	RepoPattern string
}

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

func (uc *UseCase) Execute(ctx context.Context, opts Options, content io.Reader) error {
	h := hook.Hook{
		ID:          uuid.NewString(),
		Name:        opts.Name,
		Description: opts.Description,
		Event:       hook.EventType(opts.Event),
		RepoPattern: opts.RepoPattern,
	}
	return uc.hookService.AddHook(ctx, h, content)
}