package update

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/script"
)

// UseCase is a struct that encapsulates the script service for updating scripts.
type UseCase struct {
	scriptService script.ScriptService
}

// NewUseCase creates a new instance of UseCase for updating scripts.
func NewUseCase(scriptService script.ScriptService) *UseCase {
	return &UseCase{scriptService: scriptService}
}

// Execute applies a new script identified by its ID.
func (uc *UseCase) Execute(ctx context.Context, scriptID, name string, content io.Reader) error {
	return uc.scriptService.Update(ctx, scriptID, script.Entry{
		Name:    name,
		Content: content,
	})
}
