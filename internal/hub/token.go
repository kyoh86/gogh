package hub

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/internal/cli"
	"github.com/kyoh86/xdg"
)

type TokenManager interface {
	SetGithubToken(user, token string) error
	GetGithubToken(user string) (string, error)
	DeleteGithubToken(user string) error
}

const keySeparator = "\U000F0000" // Unicode Private Use Area

type MemoryTokenManager struct {
	Host string
	m    map[string]string
}

func NewMemory(host string) (TokenManager, error) {
	return &MemoryTokenManager{Host: host, m: map[string]string{}}, nil
}

func (m *MemoryTokenManager) SetGithubToken(user, token string) error {
	m.m[m.Host+keySeparator+user] = token
	return nil
}

func (m *MemoryTokenManager) GetGithubToken(user string) (string, error) {
	return m.m[m.Host+keySeparator+user], nil
}

func (m *MemoryTokenManager) DeleteGithubToken(user string) error {
	delete(m.m, m.Host+keySeparator+user)
	return nil
}

type Keyring struct {
	ring keyring.Keyring
}

const (
	KeyringFileDir = "gogh"
)

func NewKeyring(host string) (TokenManager, error) {
	if host == "" {
		return nil, errors.New("host is empty")
	}
	serviceName := strings.Join([]string{host, env.KeyringService}, ".")

	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,

		FileDir:              filepath.Join(xdg.CacheHome(), KeyringFileDir, "keyring", host),
		FilePasswordFunc:     keyring.PromptFunc(cli.AskPassword),
		KeychainName:         serviceName,
		KeychainPasswordFunc: keyring.PromptFunc(cli.AskPassword),

		PassDir: filepath.Join(xdg.CacheHome(), KeyringFileDir, "pass", host),
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v\n", ring)
	keys, _ := ring.Keys()
	fmt.Printf("%#v\n", keys)
	return &Keyring{ring: ring}, nil
}

func (m *Keyring) SetGithubToken(user, token string) error {
	if user == "" {
		return errors.New("user is empty")
	}
	return m.ring.Set(keyring.Item{
		Key:  user,
		Data: []byte(token),
	})
}

func (m *Keyring) GetGithubToken(user string) (string, error) {
	if user == "" {
		return "", errors.New("user is empty")
	}
	envar := os.Getenv("GORDON_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	item, err := m.ring.Get(user)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", nil
		}
		return "", err
	}
	return string(item.Data), nil
}

func (m *Keyring) DeleteGithubToken(user string) error {
	return m.ring.Remove(user)
}
