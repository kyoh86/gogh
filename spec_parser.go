package gogh

import (
	"errors"
	"strings"
)

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// SpecParser will parse any string as a Spec.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewSpec" instead to build Spec.
type SpecParser struct {
	defaultServer Server
	serverMap     map[string]Server
}

// Parse a string and build a Spec.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (p *SpecParser) Parse(s string) (Spec, Server, error) {
	parts := strings.Split(s, "/")
	var name, user, host string

	var server Server
	switch len(parts) {
	case 1:
		server = p.defaultServer
		host, user, name = server.Host(), server.User(), parts[0]
	case 2:
		server = p.defaultServer
		host, user, name = server.Host(), parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
		s, ok := p.serverMap[host]
		if ok {
			server = s
		} else {
			server = Server{
				host: host,
				user: user,
			}
		}
	default:
		return Spec{}, Server{}, ErrTooManySlashes
	}
	spec, err := NewSpec(host, user, name)
	if err != nil {
		return Spec{}, Server{}, err
	}
	return spec, server, nil
}

// NewSpecParser will build Spec with a default server and alternative servers.
func NewSpecParser(defaultServer Server, servers ...Server) *SpecParser {
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
	return &SpecParser{
		serverMap:     serverMap,
		defaultServer: defaultServer,
	}
}
