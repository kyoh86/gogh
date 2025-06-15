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

	List() iter.Seq2[*Overlay, error]
	Add(ctx context.Context, ov Overlay, content io.Reader) error
	Remove(ctx context.Context, ov Overlay) error
	Open(ctx context.Context, ov Overlay) (io.ReadCloser, error)
	Set([]Overlay) error
}
