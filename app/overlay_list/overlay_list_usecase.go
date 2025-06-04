package overlay_list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase represents the overlay list use case
type UseCase struct {
	overlayStore overlay.OverlayStore
}

func NewUseCase(overlayStore overlay.OverlayStore) *UseCase {
	return &UseCase{
		overlayStore: overlayStore,
	}
}

type Overlay = overlay.Overlay

// Execute lists all overlay patterns and their files
func (uc *UseCase) Execute(ctx context.Context) iter.Seq2[*Overlay, error] {
	return uc.overlayStore.ListOverlays(ctx)
}
