package hub

import (
	"errors"
	"os"
	"strings"

	"github.com/kyoh86/gogh/env"
	keyring "github.com/zalando/go-keyring"
)

func SetGithubToken(host, user, token string) error {
	if host == "" {
		return errors.New("host is empty")
	}
	if user == "" {
		return errors.New("user is empty")
	}
	return keyring.Set(strings.Join([]string{host, env.KeyringService}, "."), user, token)
}

func GetGithubToken(host, user string) (string, error) {
	if user == "" {
		return "", errors.New("user is empty")
	}
	envar := os.Getenv("GORDON_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	return keyring.Get(strings.Join([]string{host, env.KeyringService}, "."), user)
}
