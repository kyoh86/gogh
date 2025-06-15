package hook_edit

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/hook"
)

type UseCase struct {
	hookService hook.HookService
}

func NewUseCase(hookService hook.HookService) *UseCase {
	return &UseCase{hookService: hookService}
}

// ExtractScript extracts the script of a hook by its ID and writes it to the provided writer.
func (uc *UseCase) ExtractScript(ctx context.Context, hookID string, w io.Writer) error {
	r, err := uc.hookService.Open(ctx, hookID)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// UpdateScript applies a new script to a hook identified by its ID.
func (uc *UseCase) UpdateScript(ctx context.Context, hookID string, r io.Reader) error {
	found, err := uc.hookService.Get(ctx, hookID)
	if err != nil {
		return err
	}
	found.UpdatedNow()
	return uc.hookService.Update(ctx, *found, r)
}
