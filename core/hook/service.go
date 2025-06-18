package hook

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
)

// HookService defines the hook management interface
type HookService interface {
	store.Content

	List() iter.Seq2[Hook, error]
	ListFor(reference repository.Reference, event Event) iter.Seq2[Hook, error]
	Add(ctx context.Context, entry Entry) (id string, _ error)
	Get(ctx context.Context, id string) (Hook, error)
	Update(ctx context.Context, id string, entry Entry) error
	Remove(ctx context.Context, id string) error
	Load(iter.Seq2[Hook, error]) error
}
