package filesystem_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/wfs_mock"
	"github.com/kyoh86/gogh/v4/core/workspace"
	testtarget "github.com/kyoh86/gogh/v4/infra/filesystem"
)

// TestNewOverlayService tests creating a new OverlayService
func TestNewOverlayService(t *testing.T) {
	// Create a mock WFS implementation
	mockWFS := wfs_mock.NewMockWFS()

	// Test successful creation
	_, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	// Test directory creation error
	mockWFSWithError := wfs_mock.NewMockWFS()
	mockWFSWithError.SetError("MkdirAll", "", fs.ErrPermission)

	_, err = testtarget.NewOverlayService(mockWFSWithError)
	if err == nil {
		t.Error("expected error when directory creation fails, but got nil")
	}
}

// TestAddAndListOverlays tests adding overlays and listing them
func TestAddAndListOverlays(t *testing.T) {
	mockWFS := wfs_mock.NewMockWFS()

	service, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	ctx := context.Background()

	// Test listing empty directory
	overlays, err := service.ListOverlays(ctx)
	if err != nil {
		t.Fatalf("ListOverlays failed: %v", err)
	}
	if len(overlays) != 0 {
		t.Errorf("expected empty list, got %d items", len(overlays))
	}

	// Add some overlays
	entries := []workspace.OverlayEntry{
		{Pattern: "github.com/user1/repo1", RelativePath: ".envrc"},
		{Pattern: "github.com/user2/*", RelativePath: "config/settings.json"},
	}

	for _, entry := range entries {
		content := strings.NewReader("content for " + entry.RelativePath)
		err := service.AddOverlay(ctx, entry, content)
		if err != nil {
			t.Fatalf("AddOverlay failed for %+v: %v", entry, err)
		}
	}

	// List overlays and verify
	overlays, err = service.ListOverlays(ctx)
	if err != nil {
		t.Fatalf("ListOverlays failed: %v", err)
	}

	if len(overlays) != len(entries) {
		enc, err := json.MarshalIndent(mockWFS.DirEntries(), "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal mock WFS entries: %v", err)
		}
		t.Log(string(enc))
		t.Errorf("expected %d overlays, got %d", len(entries), len(overlays))
	}

	// Check if all entries are present
	foundEntries := make(map[string]bool)
	for _, overlay := range overlays {
		key := overlay.Pattern + ":" + overlay.RelativePath
		foundEntries[key] = true
	}

	for _, entry := range entries {
		key := entry.Pattern + ":" + entry.RelativePath
		if !foundEntries[key] {
			t.Errorf("expected entry %s not found in result", key)
		}
	}

	// Test error during listing
	mockWFS.SetError("ReadDir", "", fs.ErrPermission)
	_, err = service.ListOverlays(ctx)
	if err == nil {
		t.Error("expected error during ListOverlays, but got nil")
	}
	mockWFS.SetError("ReadDir", "", nil) // Clear error
}

// TestRemoveOverlay tests removing an overlay
func TestRemoveOverlay(t *testing.T) {
	mockWFS := wfs_mock.NewMockWFS()

	service, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	ctx := context.Background()

	// Add an overlay
	entry := workspace.OverlayEntry{
		Pattern:      "github.com/user/repo",
		RelativePath: ".envrc",
	}
	err = service.AddOverlay(ctx, entry, strings.NewReader("content"))
	if err != nil {
		t.Fatalf("AddOverlay failed: %v", err)
	}

	// Remove it
	err = service.RemoveOverlay(ctx, entry)
	if err != nil {
		t.Fatalf("RemoveOverlay failed: %v", err)
	}

	// Verify it's gone
	overlays, err := service.ListOverlays(ctx)
	if err != nil {
		t.Fatalf("ListOverlays failed after removal: %v", err)
	}
	for _, overlay := range overlays {
		if overlay.Pattern == entry.Pattern && overlay.RelativePath == entry.RelativePath {
			t.Errorf("expected overlay %s:%s to be removed, but it still exists", entry.Pattern, entry.RelativePath)
			return
		}
	}

	// Test removing non-existent entry
	err = service.RemoveOverlay(ctx, workspace.OverlayEntry{Pattern: "non-existent", RelativePath: "file.txt"})
	if err == nil {
		t.Error("expected error for non-existent entry, but got nil")
	}
}

// TestEncodeDecodeFileName tests the encoding and decoding of file names
func TestEncodeDecodeFileName(t *testing.T) {
	testCases := []struct {
		name         string
		pattern      string
		forInit      bool
		relativePath string
		expected     string
	}{
		{
			name:         "simple overlay",
			pattern:      "github.com/user/repo",
			forInit:      false,
			relativePath: ".envrc",
		},
		{
			name:         "simple init",
			pattern:      "github.com/user/repo",
			forInit:      true,
			relativePath: ".envrc",
		},
		{
			name:         "with wildcard",
			pattern:      "github.com/user/*",
			forInit:      false,
			relativePath: "config.json",
		},
		{
			name:         "with special chars",
			pattern:      "github.com/user-name/repo+name",
			forInit:      true,
			relativePath: "path/with spaces/and#special$chars.txt",
		},
		{
			name:         "with unicode",
			pattern:      "github.com/user/😊",
			forInit:      false,
			relativePath: "path/to/file/日本語.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			entry := workspace.OverlayEntry{
				Pattern:      tc.pattern,
				ForInit:      tc.forInit,
				RelativePath: tc.relativePath,
			}
			encoded := testtarget.EncodeFileName(entry)

			parts := strings.SplitN(encoded, "/", 3)
			if len(parts) != 3 {
				t.Fatalf("encoded filename should have three parts: got %q", encoded)
			}
			encodedPattern, encodedType, encodedRelativePath := parts[0], parts[1], parts[2]
			if encodedPattern == "" || encodedRelativePath == "" {
				t.Fatalf("encoded filename is empty: got %q", encoded)
			}
			if encodedPattern == tc.pattern {
				t.Errorf("encoded pattern should not match original: got %q, want different", encodedPattern)
			}
			if strings.Contains(encodedPattern, "/") {
				t.Errorf("encoded pattern should not contain slashes: got %q", encodedPattern)
			}
			expectedType := "overlay"
			if tc.forInit {
				expectedType = "init"
			}
			if encodedType != expectedType {
				t.Errorf("encoded type mismatch: got %q, want %q", encodedType, expectedType)
			}
			if encodedRelativePath != tc.relativePath {
				t.Errorf("encoded relativePath mismatch: got %q, want %q", encodedRelativePath, tc.relativePath)
			}

			// Decode
			decodedEntry, err := testtarget.DecodeFileName(encoded)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}

			// Verify
			if decodedEntry.Pattern != tc.pattern {
				t.Errorf("pattern mismatch: got %q, want %q", decodedEntry.Pattern, tc.pattern)
			}
			if decodedEntry.ForInit != tc.forInit {
				t.Errorf("forInit mismatch: got %v, want %v", decodedEntry.ForInit, tc.forInit)
			}
			if decodedEntry.RelativePath != tc.relativePath {
				t.Errorf("relativePath mismatch: got %q, want %q", decodedEntry.RelativePath, tc.relativePath)
			}
		})
	}

	// Test invalid encoded filename
	_, err := testtarget.DecodeFileName("invalid-not-base64")
	if err == nil {
		t.Error("expected error for invalid encoded filename, but got nil")
	}
}

// TestFindOverlays tests the FindOverlays method
func TestFindOverlays(t *testing.T) {
	// Create a mock filesystem
	mockWFS := wfs_mock.NewMockWFS()

	service, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	ctx := context.Background()

	// Create test entries
	testEntries := []workspace.OverlayEntry{
		{
			Pattern:      "github.com/user/repo",
			ForInit:      false,
			RelativePath: "config.json",
		},
		{
			Pattern:      "github.com/user/*",
			ForInit:      false,
			RelativePath: "common.yaml",
		},
		{
			Pattern:      "github.com/org/project",
			ForInit:      true,
			RelativePath: ".envrc",
		},
		{
			Pattern:      "gitlab.com/user/*",
			ForInit:      false,
			RelativePath: "settings.json",
		},
	}

	// Add overlays to the mock filesystem
	for _, entry := range testEntries {
		encodedName := testtarget.EncodeFileName(entry)
		content := fmt.Appendf(nil, "content for %s", entry.RelativePath)
		err := mockWFS.WriteFile(encodedName, content, 0644)
		if err != nil {
			t.Fatalf("Failed to set up mock file %s: %v", encodedName, err)
		}
	}

	// Test cases
	testCases := []struct {
		name          string
		ref           repository.Reference
		expectedCount int
		expectedFiles []string
	}{
		{
			name:          "exact match",
			ref:           repository.NewReference("github.com", "user", "repo"),
			expectedCount: 2,
			expectedFiles: []string{"config.json", "common.yaml"},
		},
		{
			name:          "wildcard match",
			ref:           repository.NewReference("github.com", "user", "other-repo"),
			expectedCount: 1,
			expectedFiles: []string{"common.yaml"},
		},
		{
			name:          "init match",
			ref:           repository.NewReference("github.com", "org", "project"),
			expectedCount: 1,
			expectedFiles: []string{".envrc"},
		},
		{
			name:          "no matches",
			ref:           repository.NewReference("github.com", "different", "repo"),
			expectedCount: 0,
		},
		{
			name:          "different domain",
			ref:           repository.NewReference("gitlab.com", "user", "project"),
			expectedCount: 1,
			expectedFiles: []string{"settings.json"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test FindOverlays
			var overlays []*workspace.OverlayEntry
			for o, err := range service.FindOverlays(ctx, tc.ref) {
				if err != nil {
					t.Fatalf("FindOverlays error: %v", err)
				}
				if o == nil {
					continue
				}
				func() {
					content, err := service.OpenOverlay(ctx, *o)
					if err != nil {
						t.Fatalf("OpenOverlay failed for %s: %v", o.RelativePath, err)
					}
					// Read content
					if _, err := io.ReadAll(content); err != nil {
						t.Fatalf("Failed to read overlay content: %v", err)
					}
					overlays = append(overlays, o)
				}()
			}

			if len(overlays) != tc.expectedCount {
				t.Errorf("overlay count mismatch: got %d, want %d", len(overlays), tc.expectedCount)
			}

			// Check if expected files are in the result
			for _, expectedFile := range tc.expectedFiles {
				found := false
				for _, overlay := range overlays {
					if overlay.RelativePath == expectedFile {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected file %q not found in results", expectedFile)
				}
			}
		})
	}
}

// TestInvalidPatternHandling tests how FindOverlays handles invalid patterns
func TestInvalidPatternHandling(t *testing.T) {
	// Create a mock filesystem
	mockWFS := wfs_mock.NewMockWFS()

	service, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	ctx := context.Background()

	// Add entry with invalid pattern (this is a contrived example to force a pattern matching error)
	invalidEntry := workspace.OverlayEntry{
		Pattern:      "[invalid-pattern", // Invalid regex pattern
		ForInit:      false,
		RelativePath: "file.txt",
	}
	encodedName := testtarget.EncodeFileName(invalidEntry)
	err = mockWFS.WriteFile(encodedName, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to set up mock file: %v", err)
	}

	// Try to find overlays - shouldn't crash but might return an error
	repoRef := repository.NewReference("github.com", "user", "repo")
	if err != nil {
		t.Fatalf("Failed to parse repository reference: %v", err)
	}

	// Should handle invalid pattern errors gracefully
	foundOverlay := false
	for o, err := range service.FindOverlays(ctx, repoRef) {
		if err == nil && o != nil {
			foundOverlay = true
			break
		}
	}

	// We shouldn't find any valid overlays for this repository
	if foundOverlay {
		t.Errorf("found overlay despite invalid pattern")
	}
}
