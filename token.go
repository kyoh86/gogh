package gogh

import (
	"fmt"

	"github.com/kyoh86/gogh/v3/internal/github"
)

type Host = string
type Owner = string

type TokenManager struct {
	Hosts       Map[Host, *TokenHost] `yaml:"hosts,omitempty"`
	DefaultHost Host                  `yaml:"default_host,omitempty"`
}

type TokenHost struct {
	Owners       Map[Owner, github.Token] `yaml:"owners"`
	DefaultOwner Owner                    `yaml:"default_owner"`
}

type Map[TKey comparable, TVal any] map[TKey]TVal

func (m *Map[TKey, TVal]) Set(key TKey, val TVal) {
	if *m == nil {
		*m = map[TKey]TVal{}
	}
	(*m)[key] = val
}

func (m *Map[TKey, TVal]) Delete(key TKey) {
	if *m == nil {
		return
	}
	delete(*m, key)
}

func (m *Map[TKey, TVal]) Has(key TKey) bool {
	if *m == nil {
		return false
	}
	_, ok := (*m)[key]
	return ok
}

func (m *Map[TKey, TVal]) Get(key TKey) TVal {
	var v TVal
	return m.TryGet(key, v)
}

func (m *Map[TKey, TVal]) TryGet(key TKey, def TVal) TVal {
	if *m == nil {
		*m = map[TKey]TVal{
			key: def,
		}
		return def
	}
	if v, ok := (*m)[key]; ok {
		return v
	}
	(*m)[key] = def
	return def
}

func (t TokenManager) GetDefaultKey() (Host, Owner) {
	hostName := t.DefaultHost
	if hostName == "" {
		hostName = github.DefaultHost
	}
	host := t.Hosts.Get(hostName)
	owner := ""
	if host != nil {
		owner = host.DefaultOwner
	}
	return hostName, owner
}

func (t *TokenHost) GetDefaultToken() (Owner, github.Token) {
	if t == nil {
		return "", github.Token{}
	}
	return t.DefaultOwner, t.Owners.Get(t.DefaultOwner)
}

func (t TokenManager) Get(host, owner string) github.Token {
	return t.Hosts.TryGet(host, &TokenHost{}).Owners.Get(owner)
}

func (t *TokenManager) Set(hostName, ownerName string, token github.Token) {
	host := t.Hosts.TryGet(hostName, &TokenHost{})
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
	if !t.Hosts.Has(hostName) {
		return fmt.Errorf("host %s is not registered", hostName)
	}
	host := t.Hosts.Get(hostName)
	if !host.Owners.Has(ownerName) {
		return fmt.Errorf("owner %s is not registered in host %s", ownerName, hostName)
	}
	host.DefaultOwner = ownerName
	t.Hosts.Set(hostName, host)
	return nil
}

func (t *TokenManager) Delete(host, owner string) {
	hosts := t.Hosts.Get(host)
	if hosts == nil {
		return
	}
	hosts.Owners.Delete(owner)
	if hosts.DefaultOwner == owner {
		hosts.DefaultOwner = ""
	}
	if len(hosts.Owners) == 0 {
		t.Hosts.Delete(host)
		if t.DefaultHost == host {
			t.DefaultHost = ""
		}
	}
}

func (t TokenManager) Has(host, owner string) bool {
	hosts := t.Hosts.Get(host)
	if hosts == nil {
		return false
	}
	return hosts.Owners.Has(owner)
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
