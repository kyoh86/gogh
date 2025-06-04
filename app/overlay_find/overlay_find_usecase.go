package overlay_find

import (
	"context"
	"fmt"
	"iter"

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

type Overlay struct {
	overlay.Overlay
	Location repository.Location
}

func (uc *UseCase) Execute(ctx context.Context, refs string) iter.Seq2[*Overlay, error] {
	return func(yield func(*Overlay, error) bool) {
		refWithAlias, err := uc.referenceParser.ParseWithAlias(refs)
		if err != nil {
			yield(nil, fmt.Errorf("parsing reference '%s': %w", refs, err))
			return
		}
		ref := refWithAlias.Reference
		if refWithAlias.Alias != nil {
			ref = *refWithAlias.Alias
		}
		match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, ref)
		if err != nil {
			yield(nil, fmt.Errorf("finding repository by reference '%s': %w", refs, err))
			return
		}
		if match == nil {
			yield(nil, fmt.Errorf("repository not found for reference '%s'", refs))
			return
		}
		for overlay, err := range overlay.ForReference(uc.overlayStore.ListOverlays(ctx), ref) {
			if !yield(&Overlay{Overlay: *overlay, Location: *match}, err) {
				return
			}
		}
	}
}
