package extra_list

import (
	"context"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/app/extra_describe"
	"github.com/kyoh86/gogh/v4/core/extra"
)

// UseCase represents the extra list use case
type UseCase struct {
	extraService extra.ExtraService
	writer       io.Writer
}

// NewUseCase creates a new extra list use case
func NewUseCase(extraService extra.ExtraService, writer io.Writer) *UseCase {
	return &UseCase{
		extraService: extraService,
		writer:       writer,
	}
}

// Execute performs the extra list operation
func (uc *UseCase) Execute(ctx context.Context, asJSON bool, extraType string) error {
	var useCase interface {
		Execute(ctx context.Context, e extra_describe.Extra) error
	}

	if asJSON {
		useCase = extra_describe.NewUseCaseJSON(uc.writer)
	} else {
		useCase = extra_describe.NewUseCaseOneLine(uc.writer)
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
		if err := useCase.Execute(ctx, *e); err != nil {
			return err
		}
	}

	return nil
}
