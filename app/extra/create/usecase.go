package create

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase represents the extra create use case
type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	extraService     extra.ExtraService
	overlayService   overlay.OverlayService
	referenceParser  repository.ReferenceParser
}

// NewUsecase creates a new extra create use case
func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	extraService extra.ExtraService,
	overlayService overlay.OverlayService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
		workspaceService: workspaceService,
		finderService:    finderService,
		extraService:     extraService,
		overlayService:   overlayService,
		referenceParser:  referenceParser,
	}
}

// Options contains options for the extra create operation
type Options struct {
	Name         string
	SourceRepo   string // Repository to create from
	OverlayNames []string
}

// Execute performs the extra create operation
func (uc *Usecase) Execute(ctx context.Context, opts Options) error {
	if opts.Name == "" {
		return fmt.Errorf("name is required for named extra")
	}

	if opts.SourceRepo == "" {
		return fmt.Errorf("source repository is required")
	}

	// Parse specified repository
	sourceRef, err := uc.referenceParser.Parse(opts.SourceRepo)
	if err != nil {
		return fmt.Errorf("invalid repository reference: %w", err)
	}

	// Create extra items from overlay names
	var items []extra.Item
	for _, overlayName := range opts.OverlayNames {
		// Get overlay to validate it exists
		o, err := uc.overlayService.Get(ctx, overlayName)
		if err != nil {
			return fmt.Errorf("overlay %q not found: %w", overlayName, err)
		}

		items = append(items, extra.Item{
			OverlayID: o.ID(),
			HookID:    "", // Named extras don't have associated hooks
		})
	}

	// Create named extra
	id, err := uc.extraService.AddNamedExtra(ctx, opts.Name, *sourceRef, items)
	if err != nil {
		return fmt.Errorf("creating named extra: %w", err)
	}

	fmt.Printf("Created named extra %q with ID %s\n", opts.Name, id)
	return nil
}
