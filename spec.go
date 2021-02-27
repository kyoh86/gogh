package gogh

import (
	"net/url"
	"path"
	"regexp"
)

// Spec describes which project is in a root.
type Spec struct {
	host  string
	owner string
	name  string
}

func (s Spec) Host() string  { return s.host }
func (s Spec) Owner() string { return s.owner }
func (s Spec) Name() string  { return s.name }

func (s Spec) String() string {
	return path.Join(s.Host(), s.Owner(), s.Name())
}

var (
	ErrEmptyHost  = ErrInvalidHost("empty host")
	ErrEmptyOwner = ErrInvalidOwner("empty owner")
	ErrEmptyName  = ErrInvalidName("empty name")
)

type ErrInvalidHost string

func (e ErrInvalidHost) Error() string {
	return string(e)
}

func ValidateHost(h string) error {
	if h == "" {
		return ErrEmptyHost
	}

	u, err := url.ParseRequestURI("https://" + h)
	if err != nil {
		return ErrInvalidHost("invalid host: " + h)
	}
	if u.Host != h {
		return ErrInvalidHost("invalid host: " + h)
	}
	return nil
}

type ErrInvalidName string

func (e ErrInvalidName) Error() string {
	return string(e)
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if name == "." {
		return ErrInvalidName("'.' is reserved name")
	}
	if name == ".." {
		return ErrInvalidName("'..' is reserved name")
	}
	if invalidNameRegexp.MatchString(name) {
		return ErrInvalidName("invalid name: " + name)
	}
	return nil
}

type ErrInvalidOwner string

func (e ErrInvalidOwner) Error() string {
	return string(e)
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if owner == "" {
		return ErrEmptyOwner
	}
	if !validOwnerRegexp.MatchString(owner) {
		return ErrInvalidOwner("invalid owner: " + owner)
	}
	return nil
}

func NewSpec(host, owner, name string) (Spec, error) {
	if err := ValidateName(name); err != nil {
		return Spec{}, err
	}
	if err := ValidateOwner(owner); err != nil {
		return Spec{}, err
	}
	if err := ValidateHost(host); err != nil {
		return Spec{}, err
	}
	return Spec{
		host:  host,
		owner: owner,
		name:  name,
	}, nil
}
