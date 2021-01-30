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
	defaultServer Server
	hostTokenMap  map[string]string
}

// Parse a string and build a Description.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (d *Descriptor) Parse(s string) (Description, string, error) {
	parts := strings.Split(s, "/")
	var name, user, host string

	switch len(parts) {
	case 1:
		host, user, name = d.defaultServer.Host(), d.defaultServer.User(), parts[0]
	case 2:
		host, user, name = d.defaultServer.Host(), parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
	default:
		return Description{}, "", ErrTooManySlashes
	}
	desc, err := NewDescription(host, user, name)
	if err != nil {
		return Description{}, "", err
	}
	return desc, d.hostTokenMap[host], nil
}

// NewDescriptor will build Descriptor with a default server and alternative servers.
func NewDescriptor(defaultServer Server, servers ...Server) *Descriptor {
	h := defaultServer.Host()
	hostTokenMap := map[string]string{
		h: defaultServer.Token(),
	}
	for _, s := range servers {
		h := s.Host()
		if _, ok := hostTokenMap[h]; ok {
			continue
		}
		hostTokenMap[h] = s.Token()
	}
	return &Descriptor{
		hostTokenMap:  hostTokenMap,
		defaultServer: defaultServer,
	}
}
