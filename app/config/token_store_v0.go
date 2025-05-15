package config

import (
	"context"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/typ"
	"golang.org/x/oauth2"
)

// TokenStoreV0 is a repository for managing token configuration.
type TokenStoreV0 struct{}

type yamlTokenServiceV0 struct {
	Hosts typ.Map[string, *yamlTokenHostEntryV0] `yaml:"hosts,omitempty"`
}

type yamlTokenHostEntryV0 struct {
	Owners typ.Map[string, oauth2.Token] `yaml:"owners"`
}

// Load implements store.Loader
func (d *TokenStoreV0) Load(ctx context.Context, initial func() auth.TokenService) (auth.TokenService, error) {
	var v yamlTokenServiceV0
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
	for host, entry := range v.Hosts {
		if entry.Owners == nil {
			continue
		}
		for owner, token := range entry.Owners {
			if token.AccessToken == "" {
				continue
			}
			if err := svc.Set(host, owner, token); err != nil {
				return nil, fmt.Errorf("set token: %w", err)
			}
		}
	}
	svc.MarkSaved()
	return svc, nil
}

func NewTokenStoreV0() *TokenStoreV0 {
	return &TokenStoreV0{}
}

func (*TokenStoreV0) Source() (string, error) {
	path, err := appContextPath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

var _ store.Loader[auth.TokenService] = (*TokenStoreV0)(nil)
