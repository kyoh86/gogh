package gogh

import (
	"context"
	"strings"
)

const DefaultHost = "github.com"

type Descriptor struct {
	defaultHost string
	defaultUser string
}

func (d *Descriptor) Parse(s string) (*Description, error) {
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
		return nil, ErrTooManySlashes
	}
	if err := ValidateHost(host); err != nil {
		return nil, err
	}
	if err := ValidateUser(user); err != nil {
		return nil, err
	}
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	return &Description{
		Host: host,
		User: user,
		Name: name,
	}, nil
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
