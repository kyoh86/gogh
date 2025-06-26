package list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/hook/describe"
	"github.com/kyoh86/gogh/v4/core/hook"
)

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

func (uc *Usecase) Execute(ctx context.Context, asJSON bool) error {
	var usecase interface {
		Execute(ctx context.Context, s describe.Hook) error
	}
	if asJSON {
		usecase = describe.NewJSONUsecase(uc.writer)
	} else {
		usecase = describe.NewOnelineUsecase(uc.writer)
	}
	for s, err := range uc.hookService.List() {
		if err != nil {
			return err
		}
		if s == nil {
			continue
		}
		if err := usecase.Execute(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
