package script

import (
	"context"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/core/store"
)

// ScriptService defines the hook management interface
type ScriptService interface {
	store.Content

	List() iter.Seq2[Script, error]
	Add(ctx context.Context, entry Entry) (id string, _ error)
	Get(ctx context.Context, id string) (Script, error)
	Update(ctx context.Context, id string, entry Entry) error
	Remove(ctx context.Context, id string) error
	Open(ctx context.Context, id string) (io.ReadCloser, error)
	Load(iter.Seq2[Script, error]) error
}
