package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Content interface {
	// HasChanges returns true if the content has changes
	HasChanges() bool
	// MarkSaved marks the content as saved
	MarkSaved()
}

type Loader[T any] interface {
	Load(ctx context.Context) (T, error)
}

type Store[T Content] interface {
	Loader[T]
	Save(ctx context.Context, v T) error
}

func LoadAlternative[T Content](ctx context.Context, getDefault func() T, loaders ...Loader[T]) (T, error) {
	for i, loader := range loaders {
		svc, err := loader.Load(ctx)
		if os.IsNotExist(err) {
			continue
		}
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			var empty T
			return empty, fmt.Errorf("faield to load at %dth loader: %w", i+1, err)
		}
		return svc, nil
	}
	return getDefault(), nil
}
