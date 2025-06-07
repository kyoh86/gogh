package overlay

import (
	"context"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/core/store"
)

// OverlayService interface extends store.Content and provides overlay management logic.
type OverlayService interface {
	store.Content

	ListOverlays() iter.Seq2[*Overlay, error]
	AddOverlay(ctx context.Context, ov Overlay, content io.Reader) error
	RemoveOverlay(ctx context.Context, ov Overlay) error
	OpenOverlayContent(ctx context.Context, ov Overlay) (io.ReadCloser, error)
	SetOverlays(iter.Seq2[*Overlay, error]) error
}
