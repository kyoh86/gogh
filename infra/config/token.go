package config

import (
	"fmt"

	"github.com/kyoh86/gogh/v3/infra/github"
)

type Host = string
type Owner = string

var (
	ErrNoHost  = fmt.Errorf("no host")
	ErrNoOwner = fmt.Errorf("no owner")
)

type TokenManager struct {
	Hosts       Map[Host, *TokenHost] `yaml:"hosts,omitempty"`
	DefaultHost Host                  `yaml:"default_host,omitempty"`
}

type TokenHost struct {
	Owners       Map[Owner, github.Token] `yaml:"owners"`
	DefaultOwner Owner                    `yaml:"default_owner"`
}

func (t TokenManager) GetDefaultKey() (Host, Owner) {
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

func (t *TokenManager) GetDefaultTokenFor(hostName string) (Owner, github.Token, error) {
	if t == nil {
		return "", github.Token{}, ErrNoHost
	}
	host, ok := t.Hosts.TryGet(hostName)
	if !ok {
		return "", github.Token{}, ErrNoHost
	}
	return host.GetDefaultToken()
}

func (t *TokenHost) GetDefaultToken() (Owner, github.Token, error) {
	if t == nil {
		return "", github.Token{}, ErrNoHost
	}
	token, ok := t.Owners.TryGet(t.DefaultOwner)
	if !ok {
		return "", github.Token{}, ErrNoOwner
	}
	return t.DefaultOwner, token, nil
}

func (t TokenManager) Get(hostName, ownerName string) (github.Token, error) {
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

func (t *TokenManager) Set(hostName, ownerName string, token github.Token) {
	host := t.Hosts.GetOrSet(hostName, &TokenHost{})
	if host.DefaultOwner == "" {
		host.DefaultOwner = ownerName
	}
	host.Owners.Set(ownerName, token)
	t.Hosts.Set(hostName, host)
	if t.DefaultHost == "" {
		t.DefaultHost = hostName
	}
}

func (t *TokenManager) SetDefaultHost(hostName string) error {
	if !t.Hosts.Has(hostName) {
		return fmt.Errorf("host %s is not registered", hostName)
	}
	t.DefaultHost = hostName
	return nil
}

func (t *TokenManager) SetDefaultOwner(hostName, ownerName string) error {
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

func (t *TokenManager) Delete(hostName, ownerName string) {
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

func (t TokenManager) Has(hostName, ownerName string) bool {
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

func (t TokenManager) Entries() []TokenEntry {
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
