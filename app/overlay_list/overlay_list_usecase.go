package overlay_list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/overlay_describe"
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
		Execute(ctx context.Context, s overlay_describe.Overlay) error
	}
	if asJSON {
		if withSource {
			useCase = overlay_describe.NewUseCaseJSONWithContent(uc.overlayService, uc.writer)
		} else {
			useCase = overlay_describe.NewUseCaseJSON(uc.writer)
		}
	} else {
		if withSource {
			useCase = overlay_describe.NewUseCaseDetail(uc.overlayService, uc.writer)
		} else {
			useCase = overlay_describe.NewUseCaseOneLine(uc.writer)
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
