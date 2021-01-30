package gogh

import (
	"net/url"
	"path"
	"regexp"
)

// Spec describes which project is in a root.
type Spec struct {
	host string
	user string
	name string
}

func (s Spec) Host() string { return s.host }
func (s Spec) User() string { return s.user }
func (s Spec) Name() string { return s.name }

func (s Spec) String() string {
	return path.Join(s.Host(), s.User(), s.Name())
}

var (
	ErrEmptyHost = ErrInvalidHost("empty host")
	ErrEmptyUser = ErrInvalidUser("empty user")
	ErrEmptyName = ErrInvalidName("empty name")
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

type ErrInvalidUser string

func (e ErrInvalidUser) Error() string {
	return string(e)
}

var validUserRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateUser(user string) error {
	if user == "" {
		return ErrEmptyUser
	}
	if !validUserRegexp.MatchString(user) {
		return ErrInvalidUser("invalid user: " + user)
	}
	return nil
}

func NewSpec(host, user, name string) (Spec, error) {
	if err := ValidateName(name); err != nil {
		return Spec{}, err
	}
	if err := ValidateUser(user); err != nil {
		return Spec{}, err
	}
	if err := ValidateHost(host); err != nil {
		return Spec{}, err
	}
	return Spec{
		host: host,
		user: user,
		name: name,
	}, nil
}
