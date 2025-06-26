package remove

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/script"
)

type Usecase struct {
	scriptService script.ScriptService
}

func NewUsecase(scriptService script.ScriptService) *Usecase {
	return &Usecase{scriptService: scriptService}
}

func (uc *Usecase) Execute(ctx context.Context, scriptID string) error {
	return uc.scriptService.Remove(ctx, scriptID)
}
