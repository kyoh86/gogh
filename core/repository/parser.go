package repository

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/kyoh86/gogh/v3/util"
)

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// ReferenceParser will parse any string as a Reference.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewReference" instead to build Reference.
type ReferenceParser struct {
	defaultHost  string
	defaultOwner string
}

// ParseWithAlias parses string as a Reference and following alias.
// We can specify an alias with following '='(equal) and the alias.
//
// If it's not specified, alias will be nil value.
// If it's specified a value which equals to the ref, alias will be nil value.
func (p ReferenceParser) ParseWithAlias(s string) (*ReferenceWithAlias, error) {
	switch parts := strings.Split(s, "="); len(parts) {
	case 1:
		ref, err := p.Parse(s)
		if err != nil {
			return nil, err
		}
		return &ReferenceWithAlias{Reference: *ref}, err
	case 2:
		ref, err := p.Parse(parts[0])
		if err != nil {
			return nil, err
		}
		alias, err := ParseSiblingReference(*ref, parts[1])
		if err != nil {
			return nil, err
		}
		if alias.String() == ref.String() {
			return &ReferenceWithAlias{Reference: *ref}, err
		}
		return &ReferenceWithAlias{Reference: *ref, Alias: alias}, nil
	default:
		return nil, fmt.Errorf("invalid ref: %s", s)
	}
}

// ParseSiblingReference parses string as a repository ref and following alias
// in the same host and same owner.
func ParseSiblingReference(base Reference, s string) (*Reference, error) {
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
	return util.Ptr(NewReference(base.Host(), owner, name)), nil
}

// Parse a string and build a Reference.
//
// If the string does not have a host or a user explicitly, they will be
// replaced with a default host and default owner.
func (p ReferenceParser) Parse(s string) (*Reference, error) {
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
	return util.Ptr(NewReference(host, owner, name)), nil
}

// NewReferenceParser will build Reference with a default host and default owner.
func NewReferenceParser(defaultHost, defaultOwner string) ReferenceParser {
	return ReferenceParser{defaultHost: defaultHost, defaultOwner: defaultOwner}
}

var (
	ErrEmptyHost  = errors.New("empty host")
	ErrEmptyOwner = errors.New("empty owner")
	ErrEmptyName  = errors.New("empty name")
)

// ValidateHost validates a host string.
func ValidateHost(host string) error {
	if host == "" {
		return ErrEmptyHost
	}
	u, err := url.ParseRequestURI("https://" + host)
	if err != nil {
		return errors.New("invalid host: " + host)
	}
	if u.Host != host {
		return errors.New("invalid host: " + host)
	}
	return nil
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if owner == "" {
		return ErrEmptyOwner
	}
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("invalid owner: " + owner)
	}
	return nil
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if name == "." {
		return errors.New("'.' is reserved name")
	}
	if name == ".." {
		return errors.New("'..' is reserved name")
	}
	if invalidNameRegexp.MatchString(name) {
		return errors.New("invalid name: " + name)
	}
	return nil
}
