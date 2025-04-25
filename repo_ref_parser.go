package gogh

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// RepoRefParser will parse any string as a RepoRef.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewRepoRef" instead to build RepoRef.
type RepoRefParser struct {
	defaultHost  string
	defaultOwner string
}

// ParseWithAlias parses string as a RepoRef and following alias.
// We can specify an alias with following '='(equal) and the alias.
//
// If it's not specified, alias will be nil value.
// If it's specified a value which equals to the ref, alias will be nil value.
func (p RepoRefParser) ParseWithAlias(s string) (RepoRef, *RepoRef, error) {
	switch parts := strings.Split(s, "="); len(parts) {
	case 1:
		ref, err := p.Parse(s)
		return ref, nil, err
	case 2:
		ref, err := p.Parse(parts[0])
		if err != nil {
			return RepoRef{}, nil, err
		}
		alias, err := ParseSiblingRepoRef(ref, parts[1])
		if err != nil {
			return RepoRef{}, nil, err
		}
		if alias.String() == ref.String() {
			return ref, nil, err
		}
		return ref, &alias, nil
	default:
		return RepoRef{}, nil, fmt.Errorf("invalid ref: %s", s)
	}
}

// ParseSiblingRepoRef parses string as a repository ref and following alias
// in the same host and same owner.
func ParseSiblingRepoRef(base RepoRef, s string) (RepoRef, error) {
	parts := strings.Split(s, "/")
	var owner, name string
	switch len(parts) {
	case 1:
		owner, name = base.Owner(), parts[0]
	case 2:
		owner, name = parts[0], parts[1]
	default:
		return RepoRef{}, ErrTooManySlashes
	}
	alias, err := NewRepoRef(base.Host(), owner, name)
	if err != nil {
		return RepoRef{}, err
	}
	return alias, nil
}

// Parse a string and build a RepoRef.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default host and default owner.
func (p RepoRefParser) Parse(s string) (RepoRef, error) {
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
		return RepoRef{}, ErrTooManySlashes
	}
	ref, err := NewRepoRef(host, owner, name)
	if err != nil {
		return RepoRef{}, err
	}
	return ref, nil
}

// NewRepoRefParser will build RepoRef with a default host and default owner.
func NewRepoRefParser(defaultHost, defaultOwner string) RepoRefParser {
	return RepoRefParser{defaultHost: defaultHost, defaultOwner: defaultOwner}
}
