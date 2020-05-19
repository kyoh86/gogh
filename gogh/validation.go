package gogh

import (
	"fmt"
	"regexp"
)

type ErrorUnsupportedHost struct {
	host string
}

func (e *ErrorUnsupportedHost) Error() string {
	return fmt.Sprintf("unsupported host %q", e.host)
}

var _ error = (*ErrorUnsupportedHost)(nil)

func ValidateHost(ev Env, host string) error {
	if ev.GithubHost() != host {
		return &ErrorUnsupportedHost{host: host}
	}
	return nil
}

type ErrorInvalidName string

func (e ErrorInvalidName) Error() string {
	return string(e)
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
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

type ErrorInvalidOwner string

func (e ErrorInvalidOwner) Error() string {
	return string(e)
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return ErrorInvalidOwner("invalid owner name")
	}
	return nil
}
