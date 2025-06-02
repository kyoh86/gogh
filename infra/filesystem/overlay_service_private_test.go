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
		expected     string
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
			encoded := encodeFileName(workspace.OverlayEntry{Pattern: tc.pattern, RelativePath: tc.relativePath})

			parts := strings.SplitN(encoded, "/", 3)
			if len(parts) != 3 {
				t.Fatalf("encoded filename should have two parts: got %q", encoded)
			}
			encodedPattern, _, encodedRelativePath := parts[0], parts[1], parts[2]
			if encodedPattern == "" || encodedRelativePath == "" {
				t.Fatalf("encoded filename is empty: got %q", encoded)
			}
			if encodedPattern == tc.pattern {
				t.Errorf("encoded pattern should not match original: got %q, want different", encodedPattern)
			}
			if strings.Contains(encodedPattern, "/") {
				t.Errorf("encoded pattern should not contain slashes: got %q", encodedPattern)
			}
			if encodedRelativePath != tc.relativePath {
				t.Errorf("encoded relativePath mismatch: got %q, want %q", encodedRelativePath, tc.relativePath)
			}

			// Decode
			entry, err := decodeFileName(encoded)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}

			// Verify
			if entry.Pattern != tc.pattern {
				t.Errorf("pattern mismatch: got %q, want %q", entry.Pattern, tc.pattern)
			}
			if entry.RelativePath != tc.relativePath {
				t.Errorf("relativePath mismatch: got %q, want %q", entry.RelativePath, tc.relativePath)
			}
		})
	}

	// Test invalid encoded filename
	_, err := decodeFileName("invalid-not-base64")
	if err == nil {
		t.Error("expected error for invalid encoded filename, but got nil")
	}
}
