package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v3/typ"
)

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// ReferenceParser will parse any string as a Reference.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewReference" instead to build Reference.
type ReferenceParser interface {
	ParseWithAlias(s string) (*ReferenceWithAlias, error)
	Parse(s string) (*Reference, error)
}

type referenceParserImpl struct {
	defaultHost  string
	defaultOwner string
}

// ParseWithAlias parses string as a Reference and following alias.
// We can specify an alias with following '='(equal) and the alias.
//
// If it's not specified, alias will be nil value.
// If it's specified a value which equals to the ref, alias will be nil value.
func (p *referenceParserImpl) ParseWithAlias(s string) (*ReferenceWithAlias, error) {
	switch parts := strings.Split(s, "="); len(parts) {
	case 1:
		ref, err := p.Parse(s)
		if err != nil {
			return nil, err
		}
		return &ReferenceWithAlias{Reference: *ref}, nil
	case 2:
		ref, err := p.Parse(parts[0])
		if err != nil {
			return nil, err
		}
		alias, err := parseSiblingReference(*ref, parts[1])
		if err != nil {
			return nil, err
		}
		if alias.String() == ref.String() {
			return &ReferenceWithAlias{Reference: *ref}, nil
		}
		return &ReferenceWithAlias{Reference: *ref, Alias: alias}, nil
	default:
		return nil, fmt.Errorf("invalid ref: %s", s)
	}
}

// parseSiblingReference parses string as a repository ref and following alias
// in the same host and same owner.
func parseSiblingReference(base Reference, s string) (*Reference, error) {
	parts := strings.Split(s, "/")
	var owner, name string
	switch len(parts) {
	case 1:
		owner, name = base.Owner(), parts[0]
	case 2:
		owner, name = parts[0], parts[1]
	default:
		return nil, ErrTooManySlashes
	}
	if err := ValidateOwner(owner); err != nil {
		return nil, err
	}
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	return typ.Ptr(NewReference(base.Host(), owner, name)), nil
}

// Parse a string and build a Reference.
//
// The string will be separated host/owner/name.
// If it does not have a host or a user explicitly, they will be
// replaced with a default-host and default-owner.
func (p *referenceParserImpl) Parse(s string) (*Reference, error) {
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
		return nil, ErrTooManySlashes
	}
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	if err := ValidateOwner(owner); err != nil {
		return nil, err
	}
	if err := ValidateHost(host); err != nil {
		return nil, err
	}
	return typ.Ptr(NewReference(host, owner, name)), nil
}

// NewReferenceParser will build Reference with a default host and default owner.
func NewReferenceParser(defaultHost, defaultOwner string) ReferenceParser {
	return &referenceParserImpl{defaultHost: defaultHost, defaultOwner: defaultOwner}
}

var _ ReferenceParser = (*referenceParserImpl)(nil)
