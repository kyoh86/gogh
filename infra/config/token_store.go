package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/kyoh86/gogh/v3/infra/github"
)

type Host = string
type Owner = string

var (
	ErrNoHost  = fmt.Errorf("no host")
	ErrNoOwner = fmt.Errorf("no owner")
)

type TokenStore struct {
	Hosts       Map[Host, *TokenHostEntry] `yaml:"hosts,omitempty"`
	DefaultHost Host                       `yaml:"default_host,omitempty"`
}

type TokenHostEntry struct {
	Owners       Map[Owner, github.Token] `yaml:"owners"`
	DefaultOwner Owner                    `yaml:"default_owner"`
}

func (t TokenStore) GetDefaultKey() (Host, Owner) {
	hostName := t.DefaultHost
	if hostName == "" {
		hostName = github.DefaultHost
	}
	host, ok := t.Hosts.TryGet(hostName)
	if ok {
		return hostName, host.DefaultOwner
	}
	return hostName, ""
}

func (t *TokenStore) GetDefaultTokenFor(hostName string) (Owner, github.Token, error) {
	if t == nil {
		return "", github.Token{}, ErrNoHost
	}
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return "", github.Token{}, ErrNoHost
	}
	return host.GetDefaultToken()
}

func (t *TokenHostEntry) GetDefaultToken() (Owner, github.Token, error) {
	if t == nil {
		return "", github.Token{}, ErrNoHost
	}
	token, ok := t.Owners.TryGet(t.DefaultOwner)
	if !ok {
		return "", github.Token{}, ErrNoOwner
	}
	return t.DefaultOwner, token, nil
}

func (t TokenStore) Get(hostName, ownerName string) (github.Token, error) {
	tokenHost, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return github.Token{}, ErrNoHost
	}
	token, ok := tokenHost.Owners.TryGet(ownerName)
	if !ok {
		return github.Token{}, ErrNoOwner
	}
	return token, nil
}

func (t *TokenStore) Set(hostName, ownerName string, token github.Token) {
	host := t.Hosts.GetOrSet(hostName, &TokenHostEntry{})
	if host.DefaultOwner == "" {
		host.DefaultOwner = ownerName
	}
	host.Owners.Set(ownerName, token)
	t.Hosts.Set(hostName, host)
	if t.DefaultHost == "" {
		t.DefaultHost = hostName
	}
}

func (t *TokenStore) SetDefaultHost(hostName string) error {
	if !t.Hosts.Has(hostName) {
		return fmt.Errorf("host %s is not registered", hostName)
	}
	t.DefaultHost = hostName
	return nil
}

func (t *TokenStore) SetDefaultOwner(hostName, ownerName string) error {
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return ErrNoHost
	}
	if _, ok := host.Owners.TryGet(ownerName); !ok {
		return ErrNoOwner
	}
	host.DefaultOwner = ownerName
	t.Hosts.Set(hostName, host)
	return nil
}

func (t *TokenStore) Delete(hostName, ownerName string) {
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return
	}
	host.Owners.Delete(ownerName)
	if host.DefaultOwner == ownerName {
		host.DefaultOwner = ""
	}
	if len(host.Owners) == 0 {
		t.Hosts.Delete(hostName)
		if t.DefaultHost == hostName {
			t.DefaultHost = ""
		}
	}
}

func (t TokenStore) Has(hostName, ownerName string) bool {
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return false
	}
	return host.Owners.Has(ownerName)
}

type TokenEntry struct {
	Host  Host
	Owner Owner
	Token github.Token
}

func (e TokenEntry) String() string {
	return fmt.Sprintf("%s/%s", e.Host, e.Owner)
}

func (t TokenStore) Entries() []TokenEntry {
	var entries []TokenEntry
	for hostName, hostEntry := range t.Hosts {
		for owner, token := range hostEntry.Owners {
			entries = append(entries, TokenEntry{
				Host:  hostName,
				Owner: owner,
				Token: token,
			})
		}
	}
	return entries
}

var (
	globalTokens = TokenStore{}
	tokensOnce   sync.Once
)

func TokensPath() (string, error) {
	path, err := appContextPath("GOGH_TOKENS_PATH", os.UserCacheDir, "tokens.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

func LoadTokens() (_ *TokenStore, retErr error) {
	tokensOnce.Do(func() {
		path, err := TokensPath()
		if err != nil {
			retErr = err
			return
		}

		if err := loadYAML(path, &globalTokens); err != nil {
			retErr = err
			return
		}
	})
	return &globalTokens, retErr
}

func SaveTokens() error {
	path, err := TokensPath()
	if err != nil {
		return err
	}
	return saveYAML(path, globalTokens)
}

// GetTokenForOwner attempts to find appropriate tokens for the specified host/owner.
// Returns:
// - exactMatch: The token specifically for the requested owner (if exists)
// - candidates: Other tokens for the host that might have access to the organization
// - error: If no tokens are available or host doesn't exist
func (t TokenStore) GetTokenForOwner(hostName, ownerName string) (exactMatch *TokenEntry, candidates []TokenEntry, err error) {
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return nil, nil, ErrNoHost
	}

	// Check if we have a specific token for this owner
	if token, ok := host.Owners.TryGet(ownerName); ok {
		return &TokenEntry{
			Host:  hostName,
			Owner: ownerName,
			Token: token,
		}, nil, nil
	}

	// Gather other tokens as candidates
	for owner, token := range host.Owners {
		if owner != ownerName {
			candidates = append(candidates, TokenEntry{
				Host:  hostName,
				Owner: owner,
				Token: token,
			})
		}
	}

	if len(candidates) == 0 {
		return nil, nil, ErrNoOwner
	}

	return nil, candidates, nil
}
