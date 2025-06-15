package hook_show

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/core/hook"
)

// UseCase for running hook scripts
type UseCase struct {
	writer func(ctx context.Context, h *hook.Hook) error
}

func NewUseCase(w io.Writer, asJSON bool) *UseCase {
	if asJSON {
		enc := json.NewEncoder(w)
		return &UseCase{
			writer: func(ctx context.Context, h *hook.Hook) error {
				return enc.Encode(h)
			},
		}
	}
	return &UseCase{
		writer: func(ctx context.Context, h *hook.Hook) error {
			fmt.Printf("* [%s] %s (%s) @ %s\n", h.ID[:8], h.Name, h.Target, h.UpdatedAt.Format("2006-01-02 15:04:05"))
			return nil
		},
	}
}

func (uc *UseCase) Execute(ctx context.Context, h *hook.Hook) error {
	return uc.writer(ctx, h)
}
