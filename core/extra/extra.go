package extra

import (
	"time"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// Type represents the type of extra
type Type string

const (
	// TypeAuto represents auto-apply extra tied to a repository
	TypeAuto Type = "auto"
	// TypeNamed represents named extra used as templates
	TypeNamed Type = "named"
)

// Item represents a pair of overlay and hook
type Item struct {
	OverlayID string
	HookID    string
}

// Extra represents a collection of extra files for a repository
type Extra struct {
	id         string
	extraType  Type
	name       string                // empty for auto extra
	repository *repository.Reference // nil for named extra
	items      []Item                // overlay and hook pairs
	source     repository.Reference  // source repository
	createdAt  time.Time
}

// ID returns the unique identifier
func (e *Extra) ID() string {
	return e.id
}

// Type returns the extra type
func (e *Extra) Type() Type {
	return e.extraType
}

// Name returns the name (empty for auto extra)
func (e *Extra) Name() string {
	return e.name
}

// Repository returns the associated repository (nil for named extra)
func (e *Extra) Repository() *repository.Reference {
	if e.repository == nil {
		return nil
	}
	ref := *e.repository
	return &ref
}

// Items returns the overlay and hook pairs
func (e *Extra) Items() []Item {
	result := make([]Item, len(e.items))
	copy(result, e.items)
	return result
}

// Source returns the source repository reference
func (e *Extra) Source() repository.Reference {
	return e.source
}

// CreatedAt returns the creation time
func (e *Extra) CreatedAt() time.Time {
	return e.createdAt
}

// NewAutoExtra creates a new auto-apply extra
func NewAutoExtra(
	id string,
	repo repository.Reference,
	source repository.Reference,
	items []Item,
	createdAt time.Time,
) *Extra {
	return &Extra{
		id:         id,
		extraType:  TypeAuto,
		repository: &repo,
		items:      items,
		source:     source,
		createdAt:  createdAt,
	}
}

// NewNamedExtra creates a new named extra
func NewNamedExtra(
	id string,
	name string,
	source repository.Reference,
	items []Item,
	createdAt time.Time,
) *Extra {
	return &Extra{
		id:        id,
		extraType: TypeNamed,
		name:      name,
		items:     items,
		source:    source,
		createdAt: createdAt,
	}
}
