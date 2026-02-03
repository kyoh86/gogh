package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

type Usecase struct {
	overlayService overlay.OverlayService
}

func NewUsecase(
	overlayService overlay.OverlayService,
) *Usecase {
	return &Usecase{
		overlayService: overlayService,
	}
}

func (uc *Usecase) Execute(ctx context.Context) iter.Seq2[overlay.Overlay, error] {
	_ = ctx
	return uc.overlayService.List()
}
