package config

import (
	"fmt"
	"os"
)

// OverlayDir returns the path to the overlay directory.
func OverlayDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_PATH", os.UserConfigDir, "overlay")
	if err != nil {
		return "", fmt.Errorf("search overlay path: %w", err)
	}
	return path, nil
}
