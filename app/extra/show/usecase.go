package show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/extra/describe"
	"github.com/kyoh86/gogh/v4/core/extra"
)

// UseCase represents the extra show use case
type UseCase struct {
	extraService extra.ExtraService
	writer       io.Writer
}

// NewUseCase creates a new extra show use case
func NewUseCase(extraService extra.ExtraService, writer io.Writer) *UseCase {
	return &UseCase{
		extraService: extraService,
		writer:       writer,
	}
}

// Execute performs the extra show operation
func (uc *UseCase) Execute(ctx context.Context, identifier string, asJSON bool) error {
	// Try as ID first
	e, err := uc.extraService.Get(ctx, identifier)
	if err != nil {
		// Try as name for named extras
		e, err = uc.extraService.GetNamedExtra(ctx, identifier)
		if err != nil {
			return fmt.Errorf("extra not found: %w", err)
		}
	}

	var useCase interface {
		Execute(ctx context.Context, e describe.Extra) error
	}

	if asJSON {
		useCase = describe.NewUseCaseJSON(uc.writer)
	} else {
		useCase = describe.NewUseCaseDetail(uc.writer)
	}

	return useCase.Execute(ctx, *e)
}
