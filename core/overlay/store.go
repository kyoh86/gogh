package overlay

import (
	"context"
	"io"
)

// ContentStore is an abstraction for managing overlay content (file, DB, etc).
type ContentStore interface {
	Save(ctx context.Context, overlayID string, content io.Reader) error
	Open(ctx context.Context, overlayID string) (io.ReadCloser, error)
	Remove(ctx context.Context, overlayID string) error
}
