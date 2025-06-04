package overlay_apply

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	referenceParser  repository.ReferenceParser
	overlayStore     overlay.OverlayStore
}

func NewUseCase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	referenceParser repository.ReferenceParser,
	overlayStore overlay.OverlayStore,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
		finderService:    finderService,
		referenceParser:  referenceParser,
		overlayStore:     overlayStore,
	}
}

func (uc *UseCase) Execute(ctx context.Context, refs string, repoPattern string, forInit bool, relativePath string) error {
	refWithAlias, err := uc.referenceParser.ParseWithAlias(refs)
	if err != nil {
		return fmt.Errorf("parsing reference '%s': %w", refs, err)
	}
	match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, refWithAlias.Local())
	if err != nil {
		return fmt.Errorf("finding repository by reference '%s': %w", refs, err)
	}
	if match == nil {
		return fmt.Errorf("repository not found for reference '%s'", refs)
	}
	targetPath := filepath.Join(match.FullPath(), relativePath)

	// Open the overlay source
	source, err := uc.overlayStore.OpenOverlay(ctx, overlay.Overlay{
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
