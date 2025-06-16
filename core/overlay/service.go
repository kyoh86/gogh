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
	Add(ctx context.Context, name, relativePath string, content io.Reader) (id string, _ error)
	Remove(ctx context.Context, id string) error
	Open(ctx context.Context, id string) (io.ReadCloser, error)
	Load(iter.Seq2[*Overlay, error]) error
}
