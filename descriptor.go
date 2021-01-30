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
	serverMap     map[string]Server
}

// Parse a string and build a Description.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (d *Descriptor) Parse(s string) (Description, Server, error) {
	parts := strings.Split(s, "/")
	var name, user, host string

	var server Server
	switch len(parts) {
	case 1:
		server = d.defaultServer
		host, user, name = server.Host(), server.User(), parts[0]
	case 2:
		server = d.defaultServer
		host, user, name = server.Host(), parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
		s, ok := d.serverMap[host]
		if ok {
			server = s
		} else {
			server = Server{taggedServer{
				Host: host,
				User: user,
			}}
		}
	default:
		return Description{}, Server{}, ErrTooManySlashes
	}
	desc, err := NewDescription(host, user, name)
	if err != nil {
		return Description{}, Server{}, err
	}
	return desc, server, nil
}

// NewDescriptor will build Descriptor with a default server and alternative servers.
func NewDescriptor(defaultServer Server, servers ...Server) *Descriptor {
	h := defaultServer.Host()
	serverMap := map[string]Server{
		h: defaultServer,
	}
	for _, s := range servers {
		h := s.Host()
		if _, ok := serverMap[h]; ok {
			continue
		}
		serverMap[h] = s
	}
	return &Descriptor{
		serverMap:     serverMap,
		defaultServer: defaultServer,
	}
}
