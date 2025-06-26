package show

import (
	"context"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/app/overlay/describe"
	"github.com/kyoh86/gogh/v4/core/overlay"
)

// Usecase for running overlay overlays
type Usecase struct {
	overlayService overlay.OverlayService
	writer         io.Writer
}

func NewUsecase(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *Usecase {
	return &Usecase{
		overlayService: overlayService,
		writer:         writer,
	}
}

func (uc *Usecase) Execute(ctx context.Context, overlayID string, asJSON, withSource bool) error {
	overlay, err := uc.overlayService.Get(ctx, overlayID)
	if err != nil {
		return fmt.Errorf("get overlay by ID: %w", err)
	}
	var usecase interface {
		Execute(ctx context.Context, s describe.Overlay) error
	}
	if asJSON {
		if withSource {
			usecase = describe.NewJSONWithContentUsecase(uc.overlayService, uc.writer)
		} else {
			usecase = describe.NewJSONUsecase(uc.writer)
		}
	} else {
		if withSource {
			usecase = describe.NewDetailUsecase(uc.overlayService, uc.writer)
		} else {
			usecase = describe.NewOnelineUsecase(uc.writer)
		}
	}
	if err := usecase.Execute(ctx, overlay); err != nil {
		return fmt.Errorf("execute description: %w", err)
	}
	return nil
}
