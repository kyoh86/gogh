package gogh

import (
	"errors"
	"regexp"
)

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "." {
		return errors.New("'.' is reserved name")
	}
	if name == ".." {
		return errors.New("'..' is reserved name")
	}
	if name == "" {
		return errors.New("project name is empty")
	}
	if invalidNameRegexp.MatchString(name) {
		return errors.New("invalid project name")
	}
	return nil
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("invalid owner name")
	}
	return nil
}
