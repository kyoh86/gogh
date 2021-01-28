package gogh

import (
	"errors"
	"strings"
)

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// Descriptor will parse any string as a Description.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewDescription" instead to build Description.
type Descriptor struct {
	userHost      map[string]string
	defaultServer Server
}

// Parse a string and build a Description.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (d *Descriptor) Parse(s string) (Description, error) {
	parts := strings.Split(s, "/")
	var name, user, host string

	switch len(parts) {
	case 1:
		host, user, name = d.defaultServer.Host(), d.defaultServer.User(), parts[0]
	case 2:
		host, user, name = d.userHost[parts[0]], parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
	default:
		return Description{}, ErrTooManySlashes
	}
	return NewDescription(host, user, name)
}

// NewDescriptor will build Descriptor with a default server and alternative servers.
func NewDescriptor(defaultServer Server, servers ...Server) *Descriptor {
	u := defaultServer.User()
	userHost := map[string]string{
		u: defaultServer.Host(),
	}
	for _, s := range servers {
		u := s.User()
		if _, ok := userHost[u]; ok {
			continue
		}
		userHost[u] = s.Host()
	}
	return &Descriptor{
		userHost:      userHost,
		defaultServer: defaultServer,
	}
}
