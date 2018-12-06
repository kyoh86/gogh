package gogh

import (
	"net/url"
	"os"
	"path/filepath"
)

// gitClone git repository
func gitClone(remote *url.URL, local string, shallow bool) error {
	dir, _ := filepath.Split(local)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	args = append(args, remote.String(), local)

	return run("git", args...)
}

// gitUpdate pulls changes from remote repository
func gitUpdate(local string) error {
	return runInDir(local, "git", "pull", "--ff-only")
}
