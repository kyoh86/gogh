package config

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
)

type DefaultNameStoreV0 struct{}

type v0YAMLDefaultNameStore struct {
	Hosts map[string]struct {
		DefaultOwner string `yaml:"default_owner,omitempty"`
	} `yaml:"hosts,omitempty"`
	DefaultHost string `yaml:"default_host,omitempty"`
}

// Load implements store.Loader.
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
		return nil, fmt.Errorf("decode yaml: %w", err)
	}
	svc := initial()
	if err := svc.SetDefaultHost(v.DefaultHost); err != nil {
		return nil, err
	}
	for k, v := range v.Hosts {
		if v.DefaultOwner != "" {
			if err := svc.SetDefaultOwnerFor(k, v.DefaultOwner); err != nil {
				return nil, err
			}
		}
	}
	svc.MarkSaved()
	log.FromContext(ctx).Warnf("Default names are stored in %q which is deprecated. Please migrate to the new default names store with `gogh config migrate`.", source)
	return svc, nil
}

func (*DefaultNameStoreV0) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

func NewDefaultNameStoreV0() *DefaultNameStoreV0 {
	return &DefaultNameStoreV0{}
}

var _ store.Loader[repository.DefaultNameService] = (*DefaultNameStoreV0)(nil)
