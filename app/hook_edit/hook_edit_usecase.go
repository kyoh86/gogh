package hook_edit

import (
	"context"
	"errors"
	"io"
	"slices"

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
	var found *hook.Hook
	for h, err := range uc.hookService.ListHooks() {
		if err != nil {
			return err
		}
		if h.ID == hookID {
			found = h
			break
		}
	}
	if found == nil {
		return errors.New("hook not found")
	}
	r, err := uc.hookService.OpenHookScript(ctx, *found)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// ApplyScript applies a new script to a hook identified by its ID.
func (uc *UseCase) ApplyScript(ctx context.Context, hookID string, r io.Reader) error {
	var found *hook.Hook
	hooks := slices.Clone(collectHooks(uc.hookService))
	for _, h := range hooks {
		if h.ID == hookID {
			found = h
			break
		}
	}
	if found == nil {
		return errors.New("hook not found")
	}
	return uc.hookService.UpdateHook(ctx, *found, r)
}

func collectHooks(svc hook.HookService) []*hook.Hook {
	var out []*hook.Hook
	for h, err := range svc.ListHooks() {
		if err == nil {
			out = append(out, h)
		}
	}
	return out
}
