package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/oauth2"
)

// TokenStore is a repository for managing token configuration.
type TokenStore struct {
	filename string
}

type tomlTokenStore map[string]map[string]oauth2.Token

// Load implements auth.TokenRepository.
func (d *TokenStore) Load(ctx context.Context) (auth.TokenService, error) {
	var v tomlTokenStore
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := auth.NewTokenService()
	for host, entry := range v {
		for owner, token := range entry {
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
func (d *TokenStore) Save(ctx context.Context, ds auth.TokenService) error {
	file, err := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

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

	if err := toml.NewEncoder(file).Encode(v); err != nil {
		return err
	}
	return nil
}

func NewTokenStore(filename string) *TokenStore {
	return &TokenStore{
		filename: filename,
	}
}

func TokensPath() (string, error) {
	path, err := appContextPath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.v1.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

var _ auth.TokenStore = (*TokenStore)(nil)
