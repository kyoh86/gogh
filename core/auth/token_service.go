package auth

import (
	"errors"
	"strings"

	"github.com/kyoh86/gogh/v3/core/store"
	"golang.org/x/oauth2"
)

// Token represents an authentication token for a repository hosting service
type Token = oauth2.Token

// TokenService provides access to authentication tokens
type TokenService interface {
	// Get retrieves a token for the specified host and owner
	Get(host, owner string) (Token, error)

	// Set stores a token for the specified host and owner
	Set(host, owner string, token Token) error

	// Delete removes a token for the specified host and owner
	Delete(host, owner string) error

	// Has checks if a token exists for the specified host and owner
	Has(host, owner string) bool

	// Entries returns all stored token entries
	Entries() []TokenEntry

	store.Content
}

// TokenEntry represents a stored token with its host and owner
type TokenEntry struct {
	Host  string
	Owner string
	Token Token
}

// TokenStore is a service for saving and loading tokens
type TokenStore store.Store[TokenService]

var ErrTokenNotFound = errors.New("no token found")

type tokenServiceImpl struct {
	m       map[string]oauth2.Token
	changed bool
}

const hostOwnerSeparator = "/"

func (t tokenServiceImpl) Get(hostName, ownerName string) (oauth2.Token, error) {
	token, ok := t.m[hostName+hostOwnerSeparator+ownerName]
	if !ok {
		return oauth2.Token{}, ErrTokenNotFound
	}
	return token, nil
}

func (t *tokenServiceImpl) Set(hostName, ownerName string, token Token) error {
	t.m[hostName+hostOwnerSeparator+ownerName] = token
	t.changed = true
	return nil
}

func (t *tokenServiceImpl) Delete(hostName, ownerName string) error {
	delete(t.m, hostName+hostOwnerSeparator+ownerName)
	t.changed = true
	return nil
}

func (t tokenServiceImpl) Has(hostName, ownerName string) bool {
	_, ok := t.m[hostName+hostOwnerSeparator+ownerName]
	return ok
}

func (t tokenServiceImpl) Entries() []TokenEntry {
	entries := make([]TokenEntry, 0, len(t.m))
	for key, token := range t.m {
		words := strings.SplitN(key, hostOwnerSeparator, 2)
		if len(words) != 2 {
			continue
		}
		hostName := words[0]
		ownerName := words[1]
		entries = append(entries, TokenEntry{
			Host:  hostName,
			Owner: ownerName,
			Token: token,
		})
	}
	return entries
}

// HasChanges implements TokenService.
func (t *tokenServiceImpl) HasChanges() bool {
	return t.changed
}

// MarkSaved implements TokenService.
func (t *tokenServiceImpl) MarkSaved() {
	t.changed = false
}

func NewTokenService() TokenService {
	return &tokenServiceImpl{
		m: make(map[string]Token),
	}
}
