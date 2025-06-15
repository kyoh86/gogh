package overlay_list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase represents the overlay list use case
type UseCase struct {
	overlayService overlay.OverlayService
}

func NewUseCase(overlayService overlay.OverlayService) *UseCase {
	return &UseCase{
		overlayService: overlayService,
	}
}

type Overlay = overlay.Overlay

// Execute lists all overlay patterns and their files
func (uc *UseCase) Execute(ctx context.Context) iter.Seq2[*Overlay, error] {
	return uc.overlayService.List()
}
