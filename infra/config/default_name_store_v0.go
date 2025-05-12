package config

import (
	"context"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/store"
)

type DefaultNameStoreV0 struct {
	filename string
}

type v0YAMLDefaultNameStore struct {
	Hosts map[string]struct {
		DefaultOwner string `yaml:"default_owner,omitempty"`
	} `yaml:"hosts,omitempty"`
	DefaultHost string `yaml:"default_host,omitempty"`
}

// Load implements repository.DefaultNAmeRepositoryOld.
func (d *DefaultNameStoreV0) Load(ctx context.Context) (repository.DefaultNameService, error) {
	var v v0YAMLDefaultNameStore
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	hosts := make(map[string]string)
	for k, v := range v.Hosts {
		if v.DefaultOwner != "" {
			hosts[k] = v.DefaultOwner
		}
	}
	return &DefaultNameService{
		hosts:       hosts,
		defaultHost: v.DefaultHost,
	}, nil
}

func NewDefaultNameStoreV0(filename string) *DefaultNameStoreV0 {
	return &DefaultNameStoreV0{
		filename: filename,
	}
}

var _ store.Loader[repository.DefaultNameService] = (*DefaultNameStoreV0)(nil)
