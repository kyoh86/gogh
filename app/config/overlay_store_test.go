package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestOverlayDir(t *testing.T) {
	originalAppContextPathFunc := AppContextPathFunc
	defer func() {
		AppContextPathFunc = originalAppContextPathFunc
	}()

	testCases := []struct {
		name           string
		mockPathFunc   func(envName string, fallbackFunc func() (string, error), subdir string) (string, error)
		expectedPath   string
		expectedErrMsg string
	}{
		{
			name: "successfully returns overlay path when environment variable is set",
			mockPathFunc: func(envName string, fallbackFunc func() (string, error), subdir string) (string, error) {
				if envName != "GOGH_OVERLAY_PATH" {
					t.Errorf("expected env name to be GOGH_OVERLAY_PATH, got %s", envName)
				}
				if subdir != "overlay" {
					t.Errorf("expected subdir to be overlay, got %s", subdir)
				}
				return "/custom/path/overlay", nil
			},
			expectedPath: "/custom/path/overlay",
		},
		{
			name: "successfully returns default overlay path when environment variable is not set",
			mockPathFunc: func(envName string, fallbackFunc func() (string, error), subdir string) (string, error) {
				if envName != "GOGH_OVERLAY_PATH" {
					t.Errorf("expected env name to be GOGH_OVERLAY_PATH, got %s", envName)
				}
				if subdir != "overlay" {
					t.Errorf("expected subdir to be overlay, got %s", subdir)
				}
				// Simulate UserConfigDir returning a path
				return filepath.Join("/home/user/.config/gogh", "overlay"), nil
			},
			expectedPath: filepath.Join("/home/user/.config/gogh", "overlay"),
		},
		{
			name: "returns error when path function fails",
			mockPathFunc: func(envName string, fallbackFunc func() (string, error), subdir string) (string, error) {
				return "", errors.New("failed to get path")
			},
			expectedErrMsg: "search overlay path: failed to get path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock for AppContextPathFunc
			AppContextPathFunc = tc.mockPathFunc

			// Call the function under test
			result, err := OverlayDir()

			// Verify results
			if tc.expectedErrMsg != "" {
				if err == nil || err.Error() != tc.expectedErrMsg {
					t.Errorf("expected error %q, got %v", tc.expectedErrMsg, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if result != tc.expectedPath {
				t.Errorf("expected path %q, got %q", tc.expectedPath, result)
			}
		})
	}
}