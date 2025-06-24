package list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/hook/describe"
	"github.com/kyoh86/gogh/v4/core/hook"
)

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

func (uc *UseCase) Execute(ctx context.Context, asJSON bool) error {
	var useCase interface {
		Execute(ctx context.Context, s describe.Hook) error
	}
	if asJSON {
		useCase = describe.NewUseCaseJSON(uc.writer)
	} else {
		useCase = describe.NewUseCaseOneLine(uc.writer)
	}
	for s, err := range uc.hookService.List() {
		if err != nil {
			return err
		}
		if s == nil {
			continue
		}
		if err := useCase.Execute(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
