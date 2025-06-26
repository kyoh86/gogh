package list

import (
	"context"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/app/extra/describe"
	"github.com/kyoh86/gogh/v4/core/extra"
)

// Usecase represents the extra list use case
type Usecase struct {
	extraService extra.ExtraService
	writer       io.Writer
}

// NewUsecase creates a new extra list use case
func NewUsecase(extraService extra.ExtraService, writer io.Writer) *Usecase {
	return &Usecase{
		extraService: extraService,
		writer:       writer,
	}
}

// Execute performs the extra list operation
func (uc *Usecase) Execute(ctx context.Context, asJSON bool, extraType string) error {
	var usecase interface {
		Execute(ctx context.Context, e describe.Extra) error
	}

	if asJSON {
		usecase = describe.NewJSONUsecase(uc.writer)
	} else {
		usecase = describe.NewOnelineUsecase(uc.writer)
	}

	var list iter.Seq2[*extra.Extra, error]
	switch extraType {
	case "auto":
		list = uc.extraService.ListByType(ctx, extra.TypeAuto)
	case "named":
		list = uc.extraService.ListByType(ctx, extra.TypeNamed)
	default: // "all"
		list = uc.extraService.List(ctx)
	}

	for e, err := range list {
		if err != nil {
			return err
		}
		if e == nil {
			continue
		}
		if err := usecase.Execute(ctx, *e); err != nil {
			return err
		}
	}

	return nil
}
