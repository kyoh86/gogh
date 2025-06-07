package config

import (
	"context"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v4/core/store"
)

type FlagsStoreV0 struct{}

// Load implements store.Loader
func (d *FlagsStoreV0) Load(ctx context.Context, initial func() *Flags) (*Flags, error) {
	v := initial()
	v.RawHasChanges = true
	source, err := d.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

func NewFlagsStoreV0() *FlagsStoreV0 {
	return &FlagsStoreV0{}
}

// Source implements store.Loader
func (d *FlagsStoreV0) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_FLAG_PATH", os.UserConfigDir, "flag.yaml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

var _ store.Loader[*Flags] = (*FlagsStoreV0)(nil)
