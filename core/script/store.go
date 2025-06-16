package script

import (
	"context"
	"io"
)

// ScriptStore defines abstraction for saving, opening, and removing script source.
type ScriptStore interface {
	Save(ctx context.Context, scriptID string, content io.Reader) error
	Open(ctx context.Context, scriptID string) (io.ReadCloser, error)
	Remove(ctx context.Context, scriptID string) error
}
