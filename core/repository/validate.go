package repository

import (
	"errors"
	"net/url"
	"regexp"
)

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
