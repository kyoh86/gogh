package overlay_list

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the overlay list use case
type UseCase struct {
	overlayService workspace.OverlayService
}

func NewUseCase(overlayService workspace.OverlayService) *UseCase {
	return &UseCase{
		overlayService: overlayService,
	}
}

// Execute lists all overlay patterns and their files
func (uc *UseCase) Execute(ctx context.Context) ([]workspace.OverlayPattern, error) {
	return uc.overlayService.GetPatterns(), nil
}
