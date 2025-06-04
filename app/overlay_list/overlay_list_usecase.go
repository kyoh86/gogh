package overlay_list

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the overlay list use case
type UseCase struct {
	overlayStore workspace.OverlayStore
}

func NewUseCase(overlayStore workspace.OverlayStore) *UseCase {
	return &UseCase{
		overlayStore: overlayStore,
	}
}

type OverlayEntry = workspace.Overlay

// Execute lists all overlay patterns and their files
func (uc *UseCase) Execute(ctx context.Context) ([]OverlayEntry, error) {
	return uc.overlayStore.ListOverlays(ctx)
}
