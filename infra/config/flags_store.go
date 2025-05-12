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

// Load implements repository.DefaultNameRepositoryOld.
func (s *FlagsStore) Load(ctx context.Context) (*Flags, error) {
	v := DefaultFlags()
	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

func NewFlagsStore(filename string) *FlagsStore {
	return &FlagsStore{
		filename: filename,
	}
}

func FlagsPath() (string, error) {
	path, err := appContextPath("GOGH_FLAG_PATH", os.UserConfigDir, "flags.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

// TODO: Implement store.Store and save it is required.
var _ store.Loader[*Flags] = (*FlagsStore)(nil)
