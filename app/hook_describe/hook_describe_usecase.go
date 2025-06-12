package hook_describe

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/core/hook"
)

// UseCase for running hook scripts
type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(
	hookService hook.HookService,
) *UseCase {
	return &UseCase{
		hookService: hookService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, hookID string) (*hook.Hook, []byte, error) {
	hook, err := uc.hookService.GetHookByID(ctx, hookID)
	if err != nil {
		return nil, nil, fmt.Errorf("get hook by ID: %w", err)
	}
	src, err := uc.hookService.OpenHookScript(ctx, *hook)
	if err != nil {
		return nil, nil, fmt.Errorf("open hook script: %w", err)
	}
	defer src.Close()
	code, err := io.ReadAll(src)
	if err != nil {
		return nil, nil, fmt.Errorf("read script: %w", err)
	}
	return hook, code, nil
}
