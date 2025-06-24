package add

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

func (uc *UseCase) Execute(ctx context.Context, name string, content io.Reader) (script.Script, error) {
	e := script.Entry{
		Name:    name,
		Content: content,
	}
	id, err := uc.scriptService.Add(ctx, e)
	if err != nil {
		return nil, err
	}
	return uc.scriptService.Get(ctx, id)
}
