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

func (uc *UseCase) Execute(ctx context.Context, forInit bool, relativePath, repoPattern string) error {
	if err := uc.overlayService.RemoveOverlay(ctx, workspace.OverlayEntry{
		RepoPattern:  repoPattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	}); err != nil {
		return fmt.Errorf("removing entry %s for %s: %w", relativePath, repoPattern, err)
	}
	return nil
}
