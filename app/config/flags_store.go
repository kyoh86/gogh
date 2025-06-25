package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/store"
)

// FlagsStore is a store for flags.
type FlagsStore struct{}

// NewFlagsStore creates a new FlagsStore.
func NewFlagsStore() *FlagsStore {
	return &FlagsStore{}
}

// Load implements store.Loader
func (s *FlagsStore) Load(ctx context.Context, initial func() *Flags) (*Flags, error) {
	source, err := s.Source()
	if err != nil {
		return nil, err
	}

	v, err := loadTOMLFile[Flags](source)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Save implements repository.FlagsRepository.
func (d *FlagsStore) Save(ctx context.Context, flags *Flags, force bool) error {
	if !flags.HasChanges() && !force {
		return nil
	}
	source, err := d.Source()
	if err != nil {
		return err
	}
	if err := saveTOMLFile(source, flags); err != nil {
		return err
	}
	flags.MarkSaved()
	return nil
}

func (s *FlagsStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_FLAG_PATH", os.UserConfigDir, "flags.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

var _ store.Store[*Flags] = (*FlagsStore)(nil)
