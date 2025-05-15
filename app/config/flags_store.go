package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/pelletier/go-toml/v2"
)

// FlagsStore is a store for flags.
type FlagsStore struct{}

// Load implements store.Loader
func (s *FlagsStore) Load(ctx context.Context, initial func() *Flags) (*Flags, error) {
	v := initial()
	source, err := s.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

// NewFlagsStore creates a new FlagsStore.
func NewFlagsStore() *FlagsStore {
	return &FlagsStore{}
}

func (s *FlagsStore) Source() (string, error) {
	path, err := appContextPath("GOGH_FLAG_PATH", os.UserConfigDir, "flags.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

var _ store.Loader[*Flags] = (*FlagsStore)(nil)
