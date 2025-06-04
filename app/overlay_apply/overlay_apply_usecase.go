package overlay_apply

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	overlayStore workspace.OverlayStore
}

func NewUseCase(overlayStore workspace.OverlayStore) *UseCase {
	return &UseCase{
		overlayStore: overlayStore,
	}
}

func (uc *UseCase) Execute(ctx context.Context, repoPath string, repoPattern string, forInit bool, relativePath string) error {
	targetPath := filepath.Join(repoPath, relativePath)

	// Open the overlay source
	source, err := uc.overlayStore.OpenOverlay(ctx, workspace.Overlay{
		RepoPattern:  repoPattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	})
	if err != nil {
		return fmt.Errorf("opening overlay for repo-pattern '%s': %w", repoPattern, err)
	}
	defer source.Close()

	// Ensure the directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating directory '%s': %w", targetDir, err)
	}

	// Open the target file for writing
	target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening target file '%s': %w", targetPath, err)
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return fmt.Errorf("copying overlay content to target file '%s': %w", targetPath, err)
	}
	return nil
}
