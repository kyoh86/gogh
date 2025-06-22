package hook_show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/hook_describe"
	"github.com/kyoh86/gogh/v4/core/hook"
)

// UseCase for running hook hooks
type UseCase struct {
	hookService hook.HookService
	writer      io.Writer
}

func NewUseCase(
	hookService hook.HookService,
	writer io.Writer,
) *UseCase {
	return &UseCase{
		hookService: hookService,
		writer:      writer,
	}
}

func (uc *UseCase) Execute(ctx context.Context, hookID string, asJSON bool) error {
	hook, err := uc.hookService.Get(ctx, hookID)
	if err != nil {
		return fmt.Errorf("get hook by ID: %w", err)
	}
	var useCase interface {
		Execute(ctx context.Context, s hook_describe.Hook) error
	}
	if asJSON {
		useCase = hook_describe.NewUseCaseJSON(uc.writer)
	} else {
		useCase = hook_describe.NewUseCaseOneLine(uc.writer)
	}
	if err := useCase.Execute(ctx, hook); err != nil {
		return fmt.Errorf("execute dehookion: %w", err)
	}
	return nil
}
