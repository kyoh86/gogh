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
	servers *Servers
}

// Parse a string and build a Spec.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (p *SpecParser) Parse(s string) (Spec, Server, error) {
	parts := strings.Split(s, "/")
	var name, owner, host string

	var server Server
	switch len(parts) {
	case 1:
		s, err := p.servers.Default()
		if err != nil {
			return Spec{}, Server{}, err
		}
		server = s
		host, owner, name = server.Host(), server.User(), parts[0]
	case 2:
		s, err := p.servers.Default()
		if err != nil {
			return Spec{}, Server{}, err
		}
		server = s
		host, owner, name = server.Host(), parts[0], parts[1]
	case 3:
		host, owner, name = parts[0], parts[1], parts[2]
		s, err := p.servers.Find(host)
		if err == nil {
			server = s
		} else if errors.Is(err, ErrNoServer) || errors.Is(err, ErrServerNotFound) {
			server = Server{
				host: host,
				user: owner,
			}
		}

	default:
		return Spec{}, Server{}, ErrTooManySlashes
	}
	spec, err := NewSpec(host, owner, name)
	if err != nil {
		return Spec{}, Server{}, err
	}
	return spec, server, nil
}

// NewSpecParser will build Spec with a default server and alternative servers.
func NewSpecParser(servers *Servers) *SpecParser {
	return &SpecParser{servers: servers}
}
