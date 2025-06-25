package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
)

type DefaultNameStore struct{}

type tomlDefaultNameStore struct {
	Hosts       map[string]string `toml:"hosts,omitempty"`
	DefaultHost string            `toml:"default_host,omitempty"`
}

// Load implements repository.DefaultNameRepository.
func (d *DefaultNameStore) Load(ctx context.Context, initial func() repository.DefaultNameService) (repository.DefaultNameService, error) {
	source, err := d.Source()
	if err != nil {
		return nil, err
	}

	v, err := loadTOMLFile[tomlDefaultNameStore](source)
	if err != nil {
		return nil, err
	}

	svc := initial()
	if err := svc.SetDefaultHost(v.DefaultHost); err != nil {
		return nil, err
	}
	for host, owner := range v.Hosts {
		if err := svc.SetDefaultOwnerFor(host, owner); err != nil {
			return nil, fmt.Errorf("set default owner for %s: %w", host, err)
		}
	}
	svc.MarkSaved()
	return svc, nil
}

// Save implements repository.DefaultNameRepository.
func (d *DefaultNameStore) Save(ctx context.Context, ds repository.DefaultNameService, force bool) error {
	if !ds.HasChanges() && !force {
		return nil
	}
	source, err := d.Source()
	if err != nil {
		return err
	}
	v := tomlDefaultNameStore{
		Hosts:       ds.GetMap(),
		DefaultHost: ds.GetDefaultHost(),
	}

	if err := saveTOMLFile(source, v); err != nil {
		return err
	}
	ds.MarkSaved()
	return nil
}

func (*DefaultNameStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_DEFAULT_NAMES_PATH", os.UserConfigDir, "default_names.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search default names path: %w", err)
	}
	return path, nil
}

func NewDefaultNameStore() *DefaultNameStore {
	return &DefaultNameStore{}
}

var _ store.Store[repository.DefaultNameService] = (*DefaultNameStore)(nil)
