package overlay_add

import (
	"context"
	"fmt"
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

func (uc *UseCase) Execute(ctx context.Context, forInit bool, relativePath string, pattern string, sourceFile string) error {
	content, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("opening source file '%s': %w", sourceFile, err)
	}
	defer content.Close()
	if err := uc.overlayService.AddOverlay(ctx, workspace.OverlayEntry{
		Pattern:      pattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	}, content); err != nil {
		return fmt.Errorf("adding pattern %s: %w", pattern, err)
	}
	return nil
}
