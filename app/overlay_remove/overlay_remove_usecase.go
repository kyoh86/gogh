package overlay_remove

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	overlayService workspace.OverlayService
}

func NewUseCase(overlayService workspace.OverlayService) *UseCase {
	return &UseCase{
		overlayService: overlayService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, pattern string) error {
	if err := uc.overlayService.RemovePattern(pattern); err != nil {
		return fmt.Errorf("removing pattern %s: %w", pattern, err)
	}
	return nil
}
