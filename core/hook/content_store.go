package hook

import (
	"context"
	"io"
)

// HookScriptStore defines abstraction for saving, opening, and removing hook source scripts.
type HookScriptStore interface {
	SaveScript(ctx context.Context, h Hook, content io.Reader) (string, error)
	OpenScript(ctx context.Context, location string) (io.ReadCloser, error)
	RemoveScript(ctx context.Context, location string) error
}
