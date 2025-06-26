package show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/hook/describe"
	"github.com/kyoh86/gogh/v4/core/hook"
)

// Usecase for running hook hooks
type Usecase struct {
	hookService hook.HookService
	writer      io.Writer
}

func NewUsecase(
	hookService hook.HookService,
	writer io.Writer,
) *Usecase {
	return &Usecase{
		hookService: hookService,
		writer:      writer,
	}
}

func (uc *Usecase) Execute(ctx context.Context, hookID string, asJSON bool) error {
	hook, err := uc.hookService.Get(ctx, hookID)
	if err != nil {
		return fmt.Errorf("get hook by ID: %w", err)
	}
	var usecase interface {
		Execute(ctx context.Context, s describe.Hook) error
	}
	if asJSON {
		usecase = describe.NewJSONUsecase(uc.writer)
	} else {
		usecase = describe.NewOnelineUsecase(uc.writer)
	}
	if err := usecase.Execute(ctx, hook); err != nil {
		return fmt.Errorf("execute dehookion: %w", err)
	}
	return nil
}
