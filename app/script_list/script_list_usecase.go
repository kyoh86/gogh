package script_list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/script_describe"
	"github.com/kyoh86/gogh/v4/core/script"
)

type UseCase struct {
	scriptService script.ScriptService
	writer        io.Writer
}

func NewUseCase(
	scriptService script.ScriptService,
	writer io.Writer,
) *UseCase {
	return &UseCase{
		scriptService: scriptService,
		writer:        writer,
	}
}

func (uc *UseCase) Execute(ctx context.Context, asJSON, withSource bool) error {
	var useCase interface {
		Execute(ctx context.Context, s script_describe.Script) error
	}
	if asJSON {
		if withSource {
			useCase = script_describe.NewUseCaseJSONWithSource(uc.scriptService, uc.writer)
		} else {
			useCase = script_describe.NewUseCaseJSON(uc.writer)
		}
	} else {
		if withSource {
			useCase = script_describe.NewUseCaseDetail(uc.scriptService, uc.writer)
		} else {
			useCase = script_describe.NewUseCaseOneLine(uc.writer)
		}
	}
	for s, err := range uc.scriptService.List() {
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
