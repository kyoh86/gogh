package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/store"
	"golang.org/x/oauth2"
)

// TokenStore is a repository for managing token configuration.
type TokenStore struct{}

type tomlTokenStore map[string]map[string]oauth2.Token

func (d *TokenStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

// Load implements store.Store.
func (d *TokenStore) Load(ctx context.Context, initial func() auth.TokenService) (auth.TokenService, error) {
	source, err := d.Source()
	if err != nil {
		return nil, err
	}

	v, err := loadTOMLFile[tomlTokenStore](source)
	if err != nil {
		return nil, err
	}

	svc := initial()
	for host, entry := range *v {
		for owner, token := range entry {
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

// Save implements store.Store.
func (d *TokenStore) Save(ctx context.Context, ds auth.TokenService, force bool) error {
	if !ds.HasChanges() && !force {
		return nil
	}
	source, err := d.Source()
	if err != nil {
		return err
	}
	v := tomlTokenStore{}
	for _, entry := range ds.Entries() {
		host := entry.Host
		owner := entry.Owner
		token := entry.Token
		if _, ok := v[host]; !ok {
			v[host] = make(map[string]oauth2.Token)
		}
		v[host][owner] = token
	}

	if err := saveTOMLFile(source, v); err != nil {
		return err
	}
	ds.MarkSaved()
	return nil
}

func NewTokenStore() *TokenStore {
	return &TokenStore{}
}

var _ store.Store[auth.TokenService] = (*TokenStore)(nil)
