package apply

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase represents the extra apply use case
type Usecase struct {
	extraService     extra.ExtraService
	overlayService   overlay.OverlayService
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	referenceParser  repository.ReferenceParser
}

// NewUsecase creates a new extra apply use case
func NewUsecase(
	extraService extra.ExtraService,
	overlayService overlay.OverlayService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
		extraService:     extraService,
		overlayService:   overlayService,
		workspaceService: workspaceService,
		finderService:    finderService,
		referenceParser:  referenceParser,
	}
}

// Options contains options for the extra apply operation
type Options struct {
	Name       string // Named extra to apply
	TargetRepo string // Repository to apply to (optional, uses current directory if empty)
}

// Execute performs the extra apply operation
func (uc *Usecase) Execute(ctx context.Context, opts Options) error {
	if opts.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Get named extra
	e, err := uc.extraService.GetNamedExtra(ctx, opts.Name)
	if err != nil {
		return fmt.Errorf("getting named extra: %w", err)
	}

	// Determine target repository
	var location *repository.Location
	if opts.TargetRepo != "" {
		// Parse specified repository
		ref, err := uc.referenceParser.Parse(opts.TargetRepo)
		if err != nil {
			return fmt.Errorf("invalid repository reference: %w", err)
		}
		// Find repository location
		location, err = uc.finderService.FindByReference(ctx, uc.workspaceService, *ref)
		if err != nil {
			return fmt.Errorf("repository not found: %w", err)
		}
	} else {
		// Use current directory's repository
		cwd := "."
		location, err = uc.finderService.FindByPath(ctx, uc.workspaceService, cwd)
		if err != nil {
			return fmt.Errorf("current directory is not in a gogh-managed repository: %w", err)
		}
	}

	// Apply each overlay in the extra
	fmt.Printf("Applying extra %q to %s\n", opts.Name, location.Ref().String())

	for _, item := range e.Items() {
		// Get overlay
		o, err := uc.overlayService.Get(ctx, item.OverlayID)
		if err != nil {
			return fmt.Errorf("getting overlay %s: %w", item.OverlayID, err)
		}

		// Apply overlay
		targetPath := filepath.Join(location.FullPath(), o.RelativePath())

		// Ensure the directory exists
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return fmt.Errorf("creating directory %q: %w", targetDir, err)
		}

		// Open the overlay source
		source, err := uc.overlayService.Open(ctx, item.OverlayID)
		if err != nil {
			return fmt.Errorf("opening overlay %s: %w", item.OverlayID, err)
		}
		defer source.Close()

		// Open the target file for writing
		target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("opening target file %q: %w", targetPath, err)
		}
		defer target.Close()

		if _, err := io.Copy(target, source); err != nil {
			return fmt.Errorf("copying overlay content to target file: %w", err)
		}

		fmt.Printf("  Applied overlay %s to %s\n", o.Name(), o.RelativePath())
	}

	fmt.Printf("Successfully applied extra %q\n", opts.Name)
	return nil
}
