package filesystem

import (
	"strings"
	"testing"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// TestEncodeDecodeFileName tests the encoding and decoding of file names
func TestEncodeDecodeFileName(t *testing.T) {
	testCases := []struct {
		name         string
		pattern      string
		relativePath string
	}{
		{
			name:         "simple",
			pattern:      "github.com/user/repo",
			relativePath: ".envrc",
		},
		{
			name:         "with wildcard",
			pattern:      "github.com/user/*",
			relativePath: "config.json",
		},
		{
			name:         "with special chars",
			pattern:      "github.com/user-name/repo+name",
			relativePath: "path/with spaces/and#special$chars.txt",
		},
		{
			name:         "with unicode",
			pattern:      "github.com/user/😊",
			relativePath: "path/to/file/日本語.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			encoded := encodeFileName(tc.pattern, tc.relativePath)

			// Check encoded string doesn't contain path separator
			if strings.Contains(encoded, "/") {
				t.Errorf("encoded filename contains path separator: %s", encoded)
			}

			// Decode
			pattern, relativePath, err := decodeFileName(encoded)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}

			// Verify
			if pattern != tc.pattern {
				t.Errorf("pattern mismatch: got %q, want %q", pattern, tc.pattern)
			}
			if relativePath != tc.relativePath {
				t.Errorf("relativePath mismatch: got %q, want %q", relativePath, tc.relativePath)
			}
		})
	}

	// Test invalid encoded filename
	_, _, err := decodeFileName("invalid-not-base64")
	if err == nil {
		t.Error("expected error for invalid encoded filename, but got nil")
	}
}

// TestGetContentPath tests the getContentPath method
func TestGetContentPath(t *testing.T) {
	service := &OverlayService{}

	entry := workspace.OverlayEntry{
		Pattern:      "github.com/user/repo",
		RelativePath: ".envrc",
	}

	path := service.getContentPath(entry)

	// Verify that the filename is correctly encoded
	pattern, relativePath, err := decodeFileName(path)
	if err != nil {
		t.Fatalf("failed to decode path: %v", err)
	}

	if pattern != entry.Pattern {
		t.Errorf("pattern mismatch: got %q, want %q", pattern, entry.Pattern)
	}

	if relativePath != entry.RelativePath {
		t.Errorf("relativePath mismatch: got %q, want %q", relativePath, entry.RelativePath)
	}
}
