package show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/script/describe"
	"github.com/kyoh86/gogh/v4/core/script"
)

// Usecase for running script scripts
type Usecase struct {
	scriptService script.ScriptService
	writer        io.Writer
}

func NewUsecase(
	scriptService script.ScriptService,
	writer io.Writer,
) *Usecase {
	return &Usecase{
		scriptService: scriptService,
		writer:        writer,
	}
}

func (uc *Usecase) Execute(ctx context.Context, scriptID string, asJSON, withSource bool) error {
	script, err := uc.scriptService.Get(ctx, scriptID)
	if err != nil {
		return fmt.Errorf("get script by ID: %w", err)
	}
	var usecase interface {
		Execute(ctx context.Context, s describe.Script) error
	}
	if asJSON {
		if withSource {
			usecase = describe.NewJSONWithSourceUsecase(uc.scriptService, uc.writer)
		} else {
			usecase = describe.NewJSONUsecase(uc.writer)
		}
	} else {
		if withSource {
			usecase = describe.NewDetailUsecase(uc.scriptService, uc.writer)
		} else {
			usecase = describe.NewOnelineUsecase(uc.writer)
		}
	}
	if err := usecase.Execute(ctx, script); err != nil {
		return fmt.Errorf("execute description: %w", err)
	}
	return nil
}
