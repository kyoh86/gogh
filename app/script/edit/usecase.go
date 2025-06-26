package edit

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/script"
)

type Usecase struct {
	scriptService script.ScriptService
}

func NewUsecase(scriptService script.ScriptService) *Usecase {
	return &Usecase{scriptService: scriptService}
}

// ExtractScript extracts the script by its ID and writes it to the provided writer.
func (uc *Usecase) ExtractScript(ctx context.Context, scriptID string, w io.Writer) error {
	r, err := uc.scriptService.Open(ctx, scriptID)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// UpdateScript applies a new script identified by its ID.
func (uc *Usecase) UpdateScript(ctx context.Context, scriptID string, r io.Reader) error {
	return uc.scriptService.Update(ctx, scriptID, script.Entry{Content: r})
}
