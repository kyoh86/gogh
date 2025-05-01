package store

import (
	"context"
	"os"
)

type Store[T any] interface {
	Load(ctx context.Context) (T, error)
	Save(ctx context.Context, v T) error
}

func LoadAlternative[T any](ctx context.Context, stores ...Store[T]) (T, error) {
	for _, store := range stores {
		svc, err := store.Load(ctx)
		if os.IsNotExist(err) {
			continue
		}
		return svc, nil
	}
	var zero T
	return zero, os.ErrNotExist
}
