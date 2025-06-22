package overlay_show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/overlay_describe"
	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase for running overlay overlays
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

func (uc *UseCase) Execute(ctx context.Context, overlayID string, asJSON, withSource bool) error {
	overlay, err := uc.overlayService.Get(ctx, overlayID)
	if err != nil {
		return fmt.Errorf("get overlay by ID: %w", err)
	}
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
	if err := useCase.Execute(ctx, overlay); err != nil {
		return fmt.Errorf("execute description: %w", err)
	}
	return nil
}
