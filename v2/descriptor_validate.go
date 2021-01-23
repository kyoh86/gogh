package gogh

import (
	"errors"
	"net/url"
	"regexp"
)

var (
	ErrEmptyHost      = ErrorInvalidHost("empty description host")
	ErrEmptyUser      = ErrorInvalidUser("empty description user")
	ErrEmptyName      = ErrorInvalidName("empty description name")
	ErrTooManySlashes = errors.New("too many slashes")
)

type ErrorInvalidHost string

func (e ErrorInvalidHost) Error() string {
	return string(e)
}

func ValidateHost(h string) error {
	if h == "" {
		return ErrEmptyHost
	}

	if _, err := url.ParseRequestURI("https://" + h); err != nil {
		return ErrorInvalidHost("invalid host: " + h)
	}
	return nil
}

type ErrorInvalidName string

func (e ErrorInvalidName) Error() string {
	return string(e)
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if name == "." {
		return ErrorInvalidName("'.' is reserved name")
	}
	if name == ".." {
		return ErrorInvalidName("'..' is reserved name")
	}
	if name == "" {
		return ErrorInvalidName("project name is empty")
	}
	if invalidNameRegexp.MatchString(name) {
		return ErrorInvalidName("invalid project name")
	}
	return nil
}

type ErrorInvalidUser string

func (e ErrorInvalidUser) Error() string {
	return string(e)
}

var validUserRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateUser(user string) error {
	if user == "" {
		return ErrEmptyUser
	}
	if !validUserRegexp.MatchString(user) {
		return ErrorInvalidUser("invalid user name")
	}
	return nil
}

func ValidateDescription(host, user, name string) (*Description, error) {
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
		host: host,
		user: user,
		name: name,
	}, nil
}
