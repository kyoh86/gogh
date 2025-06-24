package edit

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/script"
)

type UseCase struct {
	scriptService script.ScriptService
}

func NewUseCase(scriptService script.ScriptService) *UseCase {
	return &UseCase{scriptService: scriptService}
}

// ExtractScript extracts the script by its ID and writes it to the provided writer.
func (uc *UseCase) ExtractScript(ctx context.Context, scriptID string, w io.Writer) error {
	r, err := uc.scriptService.Open(ctx, scriptID)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// UpdateScript applies a new script identified by its ID.
func (uc *UseCase) UpdateScript(ctx context.Context, scriptID string, r io.Reader) error {
	return uc.scriptService.Update(ctx, scriptID, script.Entry{Content: r})
}
