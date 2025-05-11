package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/pelletier/go-toml/v2"
)

type FlagsStore struct {
	filename string
}

// Load implements repository.DefaultNAmeRepositoryOld.
func (s *FlagsStore) Load(ctx context.Context) (*Flags, error) {
	var v Flags
	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Save implements repository.DefaultNAmeRepositoryOld.
func (s *FlagsStore) Save(ctx context.Context, ds *Flags) error {
	file, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(ds)
}

func NewFlagsStore(filename string) *FlagsStore {
	return &FlagsStore{
		filename: filename,
	}
}

func FlagsPath() (string, error) {
	path, err := appContextPath("GOGH_FLAG_PATH", os.UserConfigDir, AppName, "flags.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

var _ store.Store[*Flags] = (*FlagsStore)(nil)
