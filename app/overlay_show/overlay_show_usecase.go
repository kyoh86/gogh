package overlay_show

import (
	"context"
	"fmt"
	"io"
	"os"

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

func (uc *UseCase) Execute(ctx context.Context, pattern string, forInit bool, relativePath string) error {
	// Open the overlay content
	content, err := uc.overlayService.OpenOverlay(ctx, workspace.OverlayEntry{
		Pattern:      pattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	})
	if err != nil {
		return fmt.Errorf("opening overlay for pattern '%s': %w", pattern, err)
	}
	defer content.Close()
	if _, err := io.Copy(os.Stdout, content); err != nil {
		return fmt.Errorf("copying overlay content to stdout: %w", err)
	}
	return nil
}
