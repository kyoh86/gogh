package hook

import (
	"context"
	"io"
)

// HookContentStore defines abstraction for saving, opening, and removing hook source scripts.
type HookContentStore interface {
	SaveContent(ctx context.Context, h Hook, content io.Reader) (string, error)
	OpenContent(ctx context.Context, location string) (io.ReadCloser, error)
	RemoveContent(ctx context.Context, location string) error
}
