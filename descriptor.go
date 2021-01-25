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
	userHost map[string]string
	server   Server
}

// Parse a string and build a Description.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with default values which Descriptor has.
func (d *Descriptor) Parse(s string) (Description, error) {
	parts := strings.Split(s, "/")
	var name, user, host string

	switch len(parts) {
	case 1:
		host, user, name = d.server.Host(), d.server.User(), parts[0]
	case 2:
		host, user, name = d.userHost[parts[0]], parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
	default:
		return Description{}, ErrTooManySlashes
	}
	return NewDescription(host, user, name)
}

func NewDescriptor(server Server, servers ...Server) *Descriptor {
	u := server.User()
	userHost := map[string]string{
		u: server.Host(),
	}
	for _, s := range servers {
		u := s.User()
		if _, ok := userHost[u]; ok {
			continue
		}
		userHost[u] = s.Host()
	}
	return &Descriptor{
		userHost: userHost,
		server:   server,
	}
}
