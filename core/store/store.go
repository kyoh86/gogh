package store

import (
	"context"
)

type Content interface {
	// HasChanges returns true if the content has changes
	HasChanges() bool
	// MarkSaved marks the content as saved
	MarkSaved()
}

type Loader[T any] interface {
	Source() (string, error)
	Load(ctx context.Context) (T, error)
}

type Store[T Content] interface {
	Loader[T]
	Save(ctx context.Context, v T, force bool) error
}
