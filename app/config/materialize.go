package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/gogh"
)

// appContextPath returns the path to the app's configuration file.
//
// If the environment variable `envar` is set, it returns that.
// Specify a function to get the parent directory where the file will be placed, such as os.UserConfigDir.
// The `rel` is the relative path to the file from the dir.
//
// It will make the path that is formed as {getDir()}/{AppName=gogh}/{rel...}`
func appContextPath(envar string, getDir func() (string, error), rel ...string) (string, error) {
	if env := os.Getenv(envar); env != "" {
		return env, nil
	}
	dir, err := getDir()
	if err != nil {
		return "", fmt.Errorf("search app file dir for %s: %w", rel, err)
	}
	return filepath.Join(append([]string{dir, gogh.AppName}, rel...)...), nil
}
