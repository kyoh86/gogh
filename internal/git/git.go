package git

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/internal/run"
)

// Clone git repository
func Clone(remote *url.URL, local string, shallow bool) error {
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

	return run.Run("git", args...)
}

// Update pulls changes from remote repository
func Update(local string) error {
	return run.InDir(local, "git", "pull", "--ff-only")
}
