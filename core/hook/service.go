package hook

import (
	"context"
	"io"
	"iter"
)

// HookService defines the hook management interface
type HookService interface {
	List() iter.Seq2[*Hook, error]
	Add(ctx context.Context, h Hook, content io.Reader) error
	Get(ctx context.Context, id string) (*Hook, error)
	Update(ctx context.Context, h Hook, content io.Reader) error
	Remove(ctx context.Context, id string) error
	Open(ctx context.Context, id string) (io.ReadCloser, error)
	Set(iter.Seq2[*Hook, error]) error

	HasChanges() bool
	MarkSaved()
}
