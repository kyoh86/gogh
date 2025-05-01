package config

import (
	"context"
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/core/auth"
	"golang.org/x/oauth2"
)

// TokenStoreV0 is a repository for managing token configuration.
type TokenStoreV0 struct {
	filename string
}

type yamlTokenServiceV0 struct {
	Hosts Map[string, *yamlTokenHostEntryV0] `yaml:"hosts,omitempty"`
}

type yamlTokenHostEntryV0 struct {
	Owners Map[string, oauth2.Token] `yaml:"owners"`
}

// Load implements auth.TokenRepository.
func (d *TokenStoreV0) Load(ctx context.Context) (auth.TokenService, error) {
	var v yamlTokenServiceV0
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := auth.NewTokenService()
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
	return svc, nil
}

// Save implements auth.TokenRepository.
func (d *TokenStoreV0) Save(ctx context.Context, ds auth.TokenService) error {
	panic("not supported")
}

func NewTokenStoreV0(filename string) *TokenStoreV0 {
	return &TokenStoreV0{
		filename: filename,
	}
}

func TokensPathV0() (string, error) {
	path, err := appContextPath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

var _ auth.TokenStore = (*TokenStoreV0)(nil)
