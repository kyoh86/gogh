package gogh

type Token = string
type Host = string
type Owner = string

type TokenManager struct {
	Hosts       Map[Host, *TokenHost] `yaml:"hosts"`
	DefaultHost Host                  `yaml:"default_host"`
}

type TokenHost struct {
	Owners       Map[Owner, Token] `yaml:"owners"`
	DefaultOwner Owner             `yaml:"default_owner"`
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
	host := t.Hosts.Get(t.DefaultHost)
	owner := ""
	if host != nil {
		owner = host.DefaultOwner
	}
	return t.DefaultHost, owner
}

func (t *TokenHost) GetDefaultToken() (Owner, Token) {
	if t == nil {
		return "", ""
	}
	return t.DefaultOwner, t.Owners.Get(t.DefaultOwner)
}

func (t TokenManager) Get(host, owner string) Token {
	return t.Hosts.TryGet(host, &TokenHost{}).Owners.Get(owner)
}

func (t *TokenManager) Set(hostName, ownerName string, token Token) {
	host := t.Hosts.TryGet(hostName, &TokenHost{})
	if host.DefaultOwner == "" {
		host.DefaultOwner = ownerName
	}
	host.Owners.Set(ownerName, token)
	t.Hosts.Set(hostName, host)
}

func (t TokenManager) Delete(host, owner string) {
	hosts := t.Hosts.Get(host)
	if hosts == nil {
		return
	}
	hosts.Owners.Delete(owner)
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
	Token Token
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
