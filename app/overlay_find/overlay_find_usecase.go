package overlay_find

import (
	"context"
	"fmt"
	"iter"

	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// UseCase represents the create use case
type UseCase struct {
	referenceParser repository.ReferenceParser
	overlayStore    overlay.OverlayStore
}

func NewUseCase(
	referenceParser repository.ReferenceParser,
	overlayStore overlay.OverlayStore,
) *UseCase {
	return &UseCase{
		referenceParser: referenceParser,
		overlayStore:    overlayStore,
	}
}

type Overlay = overlay.Overlay

func (uc *UseCase) Execute(ctx context.Context, refs string) iter.Seq2[*Overlay, error] {
	return func(yield func(*Overlay, error) bool) {
		refWithAlias, err := uc.referenceParser.ParseWithAlias(refs)
		if err != nil {
			yield(nil, fmt.Errorf("parsing reference '%s': %w", refs, err))
			return
		}
		for overlay, err := range overlay.ForReference(uc.overlayStore.ListOverlays(ctx), refWithAlias.Local()) {
			if !yield(overlay, err) {
				return
			}
		}
	}
}
