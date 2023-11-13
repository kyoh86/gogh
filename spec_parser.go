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
	defaultHost  string
	defaultOwner string
}

// ParseWithAlias parses string as a Spec and following alias.
// We can specify an alias with following '='(equal) and the alias.
//
// If it's not specified, alias will be nil value.
// If it's specified a value which equals to the spec, alias will be nil value.
func (p SpecParser) ParseWithAlias(s string) (Spec, *Spec, error) {
	switch parts := strings.Split(s, "="); len(parts) {
	case 1:
		spec, err := p.Parse(s)
		return spec, nil, err
	case 2:
		spec, err := p.Parse(parts[0])
		if err != nil {
			return Spec{}, nil, err
		}
		alias, err := ParseSiblingSpec(spec, parts[1])
		if err != nil {
			return Spec{}, nil, err
		}
		if alias.String() == spec.String() {
			return spec, nil, err
		}
		return spec, &alias, nil
	default:
		return Spec{}, nil, fmt.Errorf("invalid spec: %s", s)
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
// replaced with a default host and default owner.
func (p SpecParser) Parse(s string) (Spec, error) {
	parts := strings.Split(s, "/")
	var host, owner, name string
	switch len(parts) {
	case 1:
		host, owner, name = p.defaultHost, p.defaultOwner, parts[0]
	case 2:
		host, owner, name = p.defaultHost, parts[0], parts[1]
	case 3:
		host, owner, name = parts[0], parts[1], parts[2]

	default:
		return Spec{}, ErrTooManySlashes
	}
	spec, err := NewSpec(host, owner, name)
	if err != nil {
		return Spec{}, err
	}
	return spec, nil
}

// NewSpecParser will build Spec with a default host and default owner.
func NewSpecParser(defaultHost, defaultOwner string) SpecParser {
	return SpecParser{defaultHost: defaultHost, defaultOwner: defaultOwner}
}
