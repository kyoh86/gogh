package overlay_add

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

func (uc *UseCase) Execute(ctx context.Context, pattern string, sourcePath string, targetPath string) error {
	patterns := uc.overlayService.GetPatterns()
	var files []workspace.OverlayFile

	// Find existing pattern
	for _, p := range patterns {
		if p.Pattern == pattern {
			files = p.Files
			break
		}
	}

	// Add new file
	files = append(files, workspace.OverlayFile{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	})

	if err := uc.overlayService.AddPattern(pattern, files); err != nil {
		return fmt.Errorf("adding pattern %s: %w", pattern, err)
	}
	return nil
}
