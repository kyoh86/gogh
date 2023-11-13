package gogh

import "fmt"

type Token string

type TokenTarget struct {
	Host  string
	Owner string
}

func (t TokenTarget) String() string {
	return t.Host + "/" + t.Owner
}

func ParseTokenTarget(s string) (TokenTarget, error) {
	var target TokenTarget
	if err := target.UnmarshalText([]byte(s)); err != nil {
		return TokenTarget{}, err
	}
	return target, nil
}

func (t TokenTarget) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *TokenTarget) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	var s string
	if err := unmarshaler(&s); err != nil {
		return err
	}
	t.Host, t.Owner = splitHostOwner(s)
	return nil
}

func splitHostOwner(s string) (host, owner string) {
	if s == "" {
		return
	}
	if s[0] == '/' {
		return s[1:], ""
	}
	for i, r := range s {
		if r == '/' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}

func (t TokenTarget) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TokenTarget) UnmarshalText(text []byte) error {
	fmt.Println(string(text))
	t.Host, t.Owner = splitHostOwner(string(text))
	fmt.Println(t)
	return nil
}

type TokenManager map[TokenTarget]Token

func (t TokenManager) Get(host, owner string) Token {
	return t[TokenTarget{Host: host, Owner: owner}]
}

func (t TokenManager) Set(host, owner string, token Token) {
	t[TokenTarget{Host: host, Owner: owner}] = token
}

func (t TokenManager) Delete(host, owner string) {
	delete(t, TokenTarget{Host: host, Owner: owner})
}

func (t TokenManager) Has(host, owner string) bool {
	_, ok := t[TokenTarget{Host: host, Owner: owner}]
	return ok
}

type TokenEntry struct {
	TokenTarget
	Token Token
}

func (t TokenManager) Entries() []TokenEntry {
	entries := make([]TokenEntry, 0, len(t))
	for k, v := range t {
		entries = append(entries, TokenEntry{
			TokenTarget: TokenTarget{
				Host:  k.Host,
				Owner: k.Owner,
			},
			Token: v,
		})
	}
	return entries
}

func (t TokenManager) MarshalYAML() (interface{}, error) {
	m := map[string]string{}
	for k, v := range t {
		m[k.String()] = string(v)
	}
	return m, nil
}

func (t *TokenManager) UnmarshalYAML(unmarshaler func(interface{}) error) error {
	m := map[string]string{}
	if err := unmarshaler(&m); err != nil {
		return err
	}
	*t = TokenManager{}
	for k, v := range m {
		var target TokenTarget
		if err := target.UnmarshalText([]byte(k)); err != nil {
			return err
		}
		fmt.Println(target)
		(*t)[target] = Token(v)
	}
	return nil
}
