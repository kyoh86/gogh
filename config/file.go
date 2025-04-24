package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppName = "gogh"

func appFilePath(envar string, getDir func() (string, error), rel ...string) (string, error) {
	if env := os.Getenv(envar); env != "" {
		return env, nil
	}
	dir, err := getDir()
	if err != nil {
		return "", fmt.Errorf("search app file dir for %s: %w", rel, err)
	}
	return filepath.Join(append([]string{dir, AppName}, rel...)...), nil
}
