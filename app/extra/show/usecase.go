package show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/extra/describe"
	"github.com/kyoh86/gogh/v4/core/extra"
)

// Usecase represents the extra show use case
type Usecase struct {
	extraService extra.ExtraService
	writer       io.Writer
}

// NewUsecase creates a new extra show use case
func NewUsecase(extraService extra.ExtraService, writer io.Writer) *Usecase {
	return &Usecase{
		extraService: extraService,
		writer:       writer,
	}
}

// Execute performs the extra show operation
func (uc *Usecase) Execute(ctx context.Context, identifier string, asJSON bool) error {
	// Try as ID first
	e, err := uc.extraService.Get(ctx, identifier)
	if err != nil {
		// Try as name for named extras
		e, err = uc.extraService.GetNamedExtra(ctx, identifier)
		if err != nil {
			return fmt.Errorf("extra not found: %w", err)
		}
	}

	var usecase interface {
		Execute(ctx context.Context, e describe.Extra) error
	}

	if asJSON {
		usecase = describe.NewJSONUsecase(uc.writer)
	} else {
		usecase = describe.NewDetailUsecase(uc.writer)
	}

	return usecase.Execute(ctx, *e)
}
