package gogh

import (
	"net/url"
	"regexp"
)

// Description describes which project is in a root.
type Description struct {
	host string
	user string
	name string
}

func (d Description) Host() string { return d.host }
func (d Description) User() string { return d.user }
func (d Description) Name() string { return d.name }

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

	if _, err := url.ParseRequestURI("https://" + h); err != nil {
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

func NewDescription(host, user, name string) (Description, error) {
	if err := ValidateHost(host); err != nil {
		return Description{}, err
	}
	if err := ValidateUser(user); err != nil {
		return Description{}, err
	}
	if err := ValidateName(name); err != nil {
		return Description{}, err
	}
	return Description{
		host: host,
		user: user,
		name: name,
	}, nil
}
