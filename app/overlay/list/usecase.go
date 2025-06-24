package list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/overlay/describe"
	"github.com/kyoh86/gogh/v4/core/overlay"
)

type UseCase struct {
	overlayService overlay.OverlayService
	writer         io.Writer
}

func NewUseCase(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *UseCase {
	return &UseCase{
		overlayService: overlayService,
		writer:         writer,
	}
}

func (uc *UseCase) Execute(ctx context.Context, asJSON, withSource bool) error {
	var useCase interface {
		Execute(ctx context.Context, s describe.Overlay) error
	}
	if asJSON {
		if withSource {
			useCase = describe.NewUseCaseJSONWithContent(uc.overlayService, uc.writer)
		} else {
			useCase = describe.NewUseCaseJSON(uc.writer)
		}
	} else {
		if withSource {
			useCase = describe.NewUseCaseDetail(uc.overlayService, uc.writer)
		} else {
			useCase = describe.NewUseCaseOneLine(uc.writer)
		}
	}
	for s, err := range uc.overlayService.List() {
		if err != nil {
			return err
		}
		if s == nil {
			continue
		}
		if err := useCase.Execute(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
