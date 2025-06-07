package overlay

import (
	"context"
	"io"
)

// ContentStore is an abstraction for managing overlay content (file, DB, etc).
type ContentStore interface {
	SaveContent(ctx context.Context, ov Overlay, content io.Reader) (string, error)
	OpenContent(ctx context.Context, location string) (io.ReadCloser, error)
	RemoveContent(ctx context.Context, location string) error
}
