package hook

import (
	"context"
	"io"
	"iter"
)

// HookService defines the hook management interface
type HookService interface {
	ListHooks() iter.Seq2[*Hook, error]
	AddHook(ctx context.Context, h Hook, content io.Reader) error
	GetHookByID(ctx context.Context, id string) (*Hook, error)
	UpdateHook(ctx context.Context, h Hook, content io.Reader) error
	RemoveHook(ctx context.Context, id string) error
	OpenHookScript(ctx context.Context, h Hook) (io.ReadCloser, error)
	SetHooks(iter.Seq2[*Hook, error]) error

	HasChanges() bool
	MarkSaved()
}
