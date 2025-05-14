package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/kyoh86/gogh/v3/core/store"
)

// LoadAlternative loads a value of type T using the provided loaders.
// It tries each loader in order until one succeeds or all fail.
// If all loaders fail, it returns the initial value provided by the initial function.
func LoadAlternative[T store.Content](ctx context.Context, initial func() T, loaders ...store.Loader[T]) (T, error) {
	for i, loader := range loaders {
		svc, err := loader.Load(ctx, initial)
		if os.IsNotExist(err) {
			continue
		}
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			var empty T
			return empty, fmt.Errorf("failed to load at %dth loader: %w", i+1, err)
		}
		return svc, nil
	}
	return initial(), nil
}
