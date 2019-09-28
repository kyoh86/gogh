package config

import (
	"os"
	"os/user"
	"path/filepath"
)

func expandPath(path string) string {
	if len(path) == 0 {
		return path
	}

	path = os.ExpandEnv(path)
	if path[0] != '~' || (len(path) > 1 && path[1] != filepath.Separator) {
		return path
	}

	user, err := user.Current()
	if err != nil {
		return path
	}

	return filepath.Join(user.HomeDir, path[1:])
}
