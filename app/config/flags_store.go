package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/pelletier/go-toml/v2"
)

// FlagsStore is a store for flags.
type FlagsStore struct{}

// NewFlagsStore creates a new FlagsStore.
func NewFlagsStore() *FlagsStore {
	return &FlagsStore{}
}

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

// Save implements repository.FlagsRepository.
func (d *FlagsStore) Save(ctx context.Context, flags *Flags, force bool) error {
	if !flags.HasChanges() && !force {
		return nil
	}
	source, err := d.Source()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(flags); err != nil {
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
