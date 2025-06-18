package script

import (
	"context"
	"io"
)

// ScriptSourceStore defines abstraction for saving, opening, and removing script source.
type ScriptSourceStore interface {
	Save(ctx context.Context, scriptID string, content io.Reader) error
	Open(ctx context.Context, scriptID string) (io.ReadCloser, error)
	Remove(ctx context.Context, scriptID string) error
}
