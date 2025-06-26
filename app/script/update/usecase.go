package update

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/script"
)

// Usecase is a struct that encapsulates the script service for updating scripts.
type Usecase struct {
	scriptService script.ScriptService
}

// NewUsecase creates a new instance of Usecase for updating scripts.
func NewUsecase(scriptService script.ScriptService) *Usecase {
	return &Usecase{scriptService: scriptService}
}

// Execute applies a new script identified by its ID.
func (uc *Usecase) Execute(ctx context.Context, scriptID, name string, content io.Reader) error {
	return uc.scriptService.Update(ctx, scriptID, script.Entry{
		Name:    name,
		Content: content,
	})
}
