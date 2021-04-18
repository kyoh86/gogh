package gogh

import (
	"errors"
	"fmt"
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

// ParseWithAlias parses string as a Spec and following alias.
// We can specify an alias with following '='(equal) and the alias.
//
// If it's not specified, alias will be nil value.
// If it's specified a value which equals to the spec, alias will be nil value.
func (p *SpecParser) ParseWithAlias(s string) (Spec, *Spec, Server, error) {
	switch parts := strings.Split(s, "="); len(parts) {
	case 1:
		spec, server, err := p.Parse(s)
		return spec, nil, server, err
	case 2:
		spec, server, err := p.Parse(parts[0])
		if err != nil {
			return Spec{}, nil, Server{}, err
		}
		alias, err := ParseSiblingSpec(spec, parts[1])
		if err != nil {
			return Spec{}, nil, Server{}, err
		}
		if alias.String() == spec.String() {
			return spec, nil, server, err
		}
		return spec, &alias, server, nil
	default:
		return Spec{}, nil, Server{}, fmt.Errorf("invalid spec: %s", s)
	}
}

// ParseSiblingSpec parses string as a repository specification
// in the same host and same owner.
func ParseSiblingSpec(base Spec, s string) (Spec, error) {
	parts := strings.Split(s, "/")
	var owner, name string
	switch len(parts) {
	case 1:
		owner, name = base.Owner(), parts[0]
	case 2:
		owner, name = parts[0], parts[1]
	default:
		return Spec{}, ErrTooManySlashes
	}
	alias, err := NewSpec(base.Host(), owner, name)
	if err != nil {
		return Spec{}, err
	}
	return alias, nil
}

// Parse a string and build a Spec.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default server.
func (p *SpecParser) Parse(s string) (Spec, Server, error) {
	parts := strings.Split(s, "/")
	var host, owner, name string
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
