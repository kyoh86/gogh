package extra

import (
	"context"
	"errors"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
)

var (
	// ErrExtraNotFound is returned when extra is not found
	ErrExtraNotFound = errors.New("extra not found")
	// ErrExtraAlreadyExists is returned when extra already exists
	ErrExtraAlreadyExists = errors.New("extra already exists")
)

// ExtraService manages extra
type ExtraService interface {
	store.Content

	// AddAutoExtra adds auto-apply extra for a repository
	AddAutoExtra(ctx context.Context, repo repository.Reference, source repository.Reference, items []Item) (string, error)

	// AddNamedExtra adds named extra
	AddNamedExtra(ctx context.Context, name string, source repository.Reference, items []Item) (string, error)

	// GetAutoExtra retrieves auto-apply extra for a repository
	GetAutoExtra(ctx context.Context, repo repository.Reference) (*Extra, error)

	// GetNamedExtra retrieves named extra by name
	GetNamedExtra(ctx context.Context, name string) (*Extra, error)

	// Get retrieves extra by ID
	Get(ctx context.Context, id string) (*Extra, error)

	// RemoveAutoExtra removes auto-apply extra for a repository
	RemoveAutoExtra(ctx context.Context, repo repository.Reference) error

	// RemoveNamedExtra removes named extra by name
	RemoveNamedExtra(ctx context.Context, name string) error

	// Remove removes extra by ID
	Remove(ctx context.Context, id string) error

	// List lists all extra
	List(ctx context.Context) iter.Seq2[*Extra, error]

	// ListByType lists extra by type
	ListByType(ctx context.Context, extraType Type) iter.Seq2[*Extra, error]

	// Load loads extras from an iterator
	Load(iter.Seq2[*Extra, error]) error
}
