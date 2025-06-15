package overlay_remove

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase represents the create use case
type UseCase struct {
	overlayService overlay.OverlayService
}

func NewUseCase(overlayService overlay.OverlayService) *UseCase {
	return &UseCase{
		overlayService: overlayService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, forInit bool, relativePath, repoPattern string) error {
	if err := uc.overlayService.Remove(ctx, overlay.Overlay{
		RepoPattern:  repoPattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	}); err != nil {
		return fmt.Errorf("removing entry %s for %s: %w", relativePath, repoPattern, err)
	}
	return nil
}
