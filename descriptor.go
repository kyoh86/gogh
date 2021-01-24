package gogh

import (
	"context"
	"errors"
	"strings"
)

const DefaultHost = "github.com"

var (
	ErrTooManySlashes = errors.New("too many slashes")
)

// Descriptor will parse any string as a Description.
//
// If it is clear that the string has host, user and name explicitly,
// use "NewDescription" instead to build Description.
type Descriptor struct {
	defaultHost string
	defaultUser string
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
		host, user, name = d.defaultHost, d.defaultUser, parts[0]
	case 2:
		host, user, name = d.defaultHost, parts[0], parts[1]
	case 3:
		host, user, name = parts[0], parts[1], parts[2]
	default:
		return Description{}, ErrTooManySlashes
	}
	return NewDescription(host, user, name)
}

func (d *Descriptor) SetDefaultUser(user string) error {
	if err := ValidateUser(user); err != nil {
		return err
	}
	d.defaultUser = user
	return nil
}

func NewDescriptor(ctx context.Context) *Descriptor {
	return &Descriptor{defaultHost: DefaultHost}
}
