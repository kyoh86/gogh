package config

import (
	"context"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/core/store"
)

type FlagsStoreV0 struct {
	filename string
}

// Load implements repository.DefaultNAmeRepositoryOld.
func (d *FlagsStoreV0) Load(ctx context.Context) (*Flags, error) {
	v := DefaultFlags()
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

// Save implements repository.DefaultNAmeRepositoryOld.
func (d *FlagsStoreV0) Save(ctx context.Context, ds *Flags) error {
	panic("not supported")
}

func NewFlagsStoreV0(filename string) *FlagsStoreV0 {
	return &FlagsStoreV0{
		filename: filename,
	}
}

func FlagsPathV0() (string, error) {
	path, err := appContextPath("GOGH_FLAG_PATH", os.UserConfigDir, AppName, "flag.yaml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

var _ store.Store[*Flags] = (*FlagsStoreV0)(nil)
