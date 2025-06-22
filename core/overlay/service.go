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

	List() iter.Seq2[Overlay, error]
	Add(ctx context.Context, entry Entry) (id string, _ error)
	Get(ctx context.Context, idlike string) (Overlay, error)
	Update(ctx context.Context, idlike string, entry Entry) error
	Remove(ctx context.Context, idlike string) error
	Open(ctx context.Context, idlike string) (io.ReadCloser, error)
	Load(iter.Seq2[Overlay, error]) error
}
