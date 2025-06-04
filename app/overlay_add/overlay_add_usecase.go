package overlay_add

import (
	"context"
	"fmt"
	"os"

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

func (uc *UseCase) Execute(ctx context.Context, forInit bool, relativePath string, repoPattern string, sourceFile string) error {
	content, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("opening source file '%s': %w", sourceFile, err)
	}
	defer content.Close()
	if err := uc.overlayStore.AddOverlay(ctx, workspace.Overlay{
		RepoPattern:  repoPattern,
		ForInit:      forInit,
		RelativePath: relativePath,
	}, content); err != nil {
		return fmt.Errorf("adding repo-pattern %s: %w", repoPattern, err)
	}
	return nil
}
