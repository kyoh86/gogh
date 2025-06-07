package hook_apply

import (
	"context"
	"fmt"
	"io"

	lua "github.com/Shopify/go-lua"
	"github.com/kyoh86/gogh/v4/core/hook"
)

// UseCase for running hook scripts
type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

func (uc *UseCase) Execute(ctx context.Context, h hook.Hook, env map[string]string) error {
	src, err := uc.hookService.OpenHookScript(ctx, h)
	if err != nil {
		return fmt.Errorf("open hook script: %w", err)
	}
	defer src.Close()
	code, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}
	l := lua.NewState()
	lua.OpenLibraries(l)
	// 必要ならGo関数をLuaにExpose（例: リポジトリパス、環境変数、ログ出力など）
	// e.g., l.Register("gogh_log", ...)
	if err := lua.DoString(l, string(code)); err != nil {
		return fmt.Errorf("run Lua: %w", err)
	}
	return nil
}

