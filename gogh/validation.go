package gogh

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "." || name == ".." {
		return errors.New("'.' or '..' is reserved name")
	}
	if name == "" {
		return errors.New("empty project name")
	}
	if invalidNameRegexp.MatchString(name) {
		return errors.New("project name may only contain alphanumeric characters, dots or hyphens")
	}
	return nil
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen")
	}
	return nil
}

func ValidateRoot(root []string) error {
	for i, v := range root {
		path := filepath.Clean(v)
		_, err := os.Stat(path)
		switch {
		case err == nil:
			root[i], err = filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
		case os.IsNotExist(err):
			root[i] = path
		default:
			return err
		}
	}
	if len(root) == 0 {
		return errors.New("no root")
	}

	return nil
}

func ValidateContext(ctx Context) error {
	if err := ValidateRoot(ctx.Root()); err != nil {
		return err
	}
	if err := ValidateOwner(ctx.GitHubUser()); err != nil {
		return err
	}
	return nil
}
