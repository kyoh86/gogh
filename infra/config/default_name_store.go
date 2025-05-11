package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/pelletier/go-toml/v2"
)

type DefaultNameStore struct {
	filename string
}

type tomlDefaultNameStore struct {
	Hosts       map[string]string `toml:"hosts,omitempty"`
	DefaultHost string            `toml:"default_host,omitempty"`
}

// Load implements repository.DefaultNameRepository.
func (d *DefaultNameStore) Load(ctx context.Context) (repository.DefaultNameService, error) {
	var v tomlDefaultNameStore
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	return &DefaultNameService{
		hosts:       v.Hosts,
		defaultHost: v.DefaultHost,
	}, nil
}

// Save implements repository.DefaultNameRepository.
func (d *DefaultNameStore) Save(ctx context.Context, ds repository.DefaultNameService) error {
	file, err := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	v := tomlDefaultNameStore{
		Hosts:       ds.GetMap(),
		DefaultHost: ds.GetDefaultHost(),
	}

	if err := toml.NewEncoder(file).Encode(v); err != nil {
		return err
	}
	return nil
}

func DefaultNamesPath() (string, error) {
	path, err := appContextPath("GOGH_DEFAULT_NAMES_PATH", os.UserConfigDir, AppName, "default_names.v4.yaml")
	if err != nil {
		return "", fmt.Errorf("search default names path: %w", err)
	}
	return path, nil
}

func NewDefaultNameStore(filename string) *DefaultNameStore {
	return &DefaultNameStore{
		filename: filename,
	}
}

var _ repository.DefaultNameStore = (*DefaultNameStore)(nil)
