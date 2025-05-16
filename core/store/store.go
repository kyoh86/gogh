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
	Load(ctx context.Context, initial func() T) (T, error)
}

type Saver[T Content] interface {
	Source() (string, error)
	Save(ctx context.Context, v T, force bool) error
}

type Store[T Content] interface {
	Loader[T]
	Saver[T]
}
