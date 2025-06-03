package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReplaceTildeWithHome replaces a tilde (~) at the beginning of a path with the user's home directory.
func ReplaceTildeWithHome(p string) (string, error) {
	runes := []rune(p)
	switch len(runes) {
	case 0:
		return p, nil
	case 1:
		if runes[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("searching user home dir: %w", err)
			}
			return homeDir, nil
		}
	default:
		if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("searching user home dir: %w", err)
			}
			return filepath.Join(homeDir, string(runes[2:])), nil
		}
	}
	return p, nil
}
