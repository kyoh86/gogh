package remove

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// Usecase represents the create use case
type Usecase struct {
	overlayService overlay.OverlayService
}

func NewUsecase(overlayService overlay.OverlayService) *Usecase {
	return &Usecase{
		overlayService: overlayService,
	}
}

func (uc *Usecase) Execute(ctx context.Context, overlayID string) error {
	return uc.overlayService.Remove(ctx, overlayID)
}
