package gogh

import (
	"errors"
	"os"
	"path/filepath"
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

func ValidateRoot(root string) (string, error) {
	path := filepath.Clean(root)
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return filepath.EvalSymlinks(path)
	case os.IsNotExist(err):
		return path, nil
	default:
		return "", err
	}
}

func ValidateRoots(roots []string) error {
	for i, v := range roots {
		r, err := ValidateRoot(v)
		if err != nil {
			return err
		}
		roots[i] = r
	}
	if len(roots) == 0 {
		return errors.New("no root")
	}

	return nil
}
