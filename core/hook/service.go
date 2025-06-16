package hook

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/store"
)

// HookService defines the hook management interface
type HookService interface {
	store.Content

	List() iter.Seq2[*Hook, error]
	Add(
		ctx context.Context,
		name string,
		repoPattern string,
		triggerEvent Event,
		operationType OperationType,
		operationID string,
	) (id string, _ error)
	Get(ctx context.Context, id string) (*Hook, error)
	Update(
		ctx context.Context,
		id string,
		name string,
		repoPattern string,
		triggerEvent Event,
		operationType OperationType,
		operationID string,
	) error
	Remove(ctx context.Context, id string) error
	Load(iter.Seq2[*Hook, error]) error
}
