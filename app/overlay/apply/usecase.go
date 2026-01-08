package apply

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase represents the create use case
type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	referenceParser  repository.ReferenceParser
	overlayService   overlay.OverlayService
}

func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	referenceParser repository.ReferenceParser,
	overlayService overlay.OverlayService,
) *Usecase {
	return &Usecase{
		workspaceService: workspaceService,
		finderService:    finderService,
		referenceParser:  referenceParser,
		overlayService:   overlayService,
	}
}

func (uc *Usecase) Execute(ctx context.Context, refStr string, overlayID string) error {
	refWithAlias, err := uc.referenceParser.ParseWithAlias(refStr)
	if err != nil {
		return fmt.Errorf("parsing reference '%s': %w", refStr, err)
	}
	match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, refWithAlias.Local())
	if err != nil {
		return fmt.Errorf("finding repository by reference '%s': %w", refWithAlias.Local().String(), err)
	}
	return uc.Apply(ctx, match, overlayID)
}

func (uc *Usecase) Apply(ctx context.Context, location *repository.Location, overlayID string) error {
	if location == nil {
		return errors.New("repository not found")
	}

	overlay, err := uc.overlayService.Get(ctx, overlayID)
	if err != nil {
		return fmt.Errorf("getting overlay with ID '%s': %w", overlayID, err)
	}

	targetPath := filepath.Join(location.FullPath(), overlay.RelativePath())

	// Open the overlay source
	source, err := uc.overlayService.Open(ctx, overlayID)
	if err != nil {
		return fmt.Errorf("opening overlay with ID '%s': %w", overlayID, err)
	}
	defer source.Close()

	// Ensure the directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("creating directory '%s': %w", targetDir, err)
	}

	// Open the target file for writing
	target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("opening target file '%s': %w", targetPath, err)
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return fmt.Errorf("copying overlay content to target file '%s': %w", targetPath, err)
	}
	return nil
}
