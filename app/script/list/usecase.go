package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/script"
)

type Usecase struct {
	scriptService script.ScriptService
}

func NewUsecase(
	scriptService script.ScriptService,
) *Usecase {
	return &Usecase{
		scriptService: scriptService,
	}
}

func (uc *Usecase) Execute(ctx context.Context) iter.Seq2[script.Script, error] {
	_ = ctx
	return uc.scriptService.List()
}
