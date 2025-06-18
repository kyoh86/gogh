package script_show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/script_describe"
	"github.com/kyoh86/gogh/v4/core/script"
)

// UseCase for running script scripts
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

func (uc *UseCase) Execute(ctx context.Context, scriptID string, asJSON, withSource bool) error {
	script, err := uc.scriptService.Get(ctx, scriptID)
	if err != nil {
		return fmt.Errorf("get script by ID: %w", err)
	}
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
	if err := useCase.Execute(ctx, script); err != nil {
		return fmt.Errorf("execute description: %w", err)
	}
	return nil
}
