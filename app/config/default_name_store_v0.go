package config

import (
	"context"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/store"
)

type DefaultNameStoreV0 struct{}

type v0YAMLDefaultNameStore struct {
	Hosts map[string]struct {
		DefaultOwner string `yaml:"default_owner,omitempty"`
	} `yaml:"hosts,omitempty"`
	DefaultHost string `yaml:"default_host,omitempty"`
}

// Load implements repository.DefaultNAmeRepositoryOld.
func (d *DefaultNameStoreV0) Load(ctx context.Context, initial func() repository.DefaultNameService) (repository.DefaultNameService, error) {
	var v v0YAMLDefaultNameStore
	source, err := d.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := initial()
	svc.SetDefaultHost(v.DefaultHost)
	for k, v := range v.Hosts {
		if v.DefaultOwner != "" {
			svc.SetDefaultOwnerFor(k, v.DefaultOwner)
		}
	}
	svc.MarkSaved()
	return svc, nil
}

func (*DefaultNameStoreV0) Source() (string, error) {
	path, err := appContextPath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

func NewDefaultNameStoreV0() *DefaultNameStoreV0 {
	return &DefaultNameStoreV0{}
}

var _ store.Loader[repository.DefaultNameService] = (*DefaultNameStoreV0)(nil)
