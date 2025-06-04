package overlay_show

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase represents the create use case
type UseCase struct {
	overlayStore overlay.OverlayStore
}

func NewUseCase(overlayStore overlay.OverlayStore) *UseCase {
	return &UseCase{
		overlayStore: overlayStore,
	}
}

func (uc *UseCase) Execute(ctx context.Context, repoPattern string, forInit bool, relativePath string) error {
	// Open the overlay content
	content, err := uc.overlayStore.OpenOverlay(ctx, overlay.Overlay{
		RepoPattern:  repoPattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	})
	if err != nil {
		return fmt.Errorf("opening overlay for repo-pattern '%s': %w", repoPattern, err)
	}
	defer content.Close()
	if _, err := io.Copy(os.Stdout, content); err != nil {
		return fmt.Errorf("copying overlay content to stdout: %w", err)
	}
	return nil
}
