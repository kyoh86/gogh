package script_remove

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/script"
)

type UseCase struct {
	scriptService script.ScriptService
}

func NewUseCase(scriptService script.ScriptService) *UseCase {
	return &UseCase{scriptService: scriptService}
}

func (uc *UseCase) Execute(ctx context.Context, scriptID string) error {
	return uc.scriptService.Remove(ctx, scriptID)
}
