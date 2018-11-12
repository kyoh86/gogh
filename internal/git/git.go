package git

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/internal/execute"
)

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

	return execute.Run("git", args...)
}
func Update(local string) error {
	return execute.RunInDir(local, "git", "pull", "--ff-only")
}
