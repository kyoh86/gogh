package filesystem_test

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

// TestGetOverlayContent tests getting overlay content
func TestGetOverlayContent(t *testing.T) {
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
	expectedContent := "GOPATH=/go\nHOME=/home/user"
	err = service.AddOverlay(ctx, entry, strings.NewReader(expectedContent))
	if err != nil {
		t.Fatalf("AddOverlay failed: %v", err)
	}

	// Get content
	content, err := service.GetOverlayContent(ctx, entry)
	if err != nil {
		t.Fatalf("GetOverlayContent failed: %v", err)
	}
	defer content.Close()

	data, err := io.ReadAll(content)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}

	if string(data) != expectedContent {
		t.Errorf("content mismatch: got %q, want %q", string(data), expectedContent)
	}

	// Test error
	nonExistentEntry := workspace.OverlayEntry{
		Pattern:      "non-existent",
		RelativePath: "file.txt",
	}
	_, err = service.GetOverlayContent(ctx, nonExistentEntry)
	if err == nil {
		t.Error("expected error for non-existent entry, but got nil")
	}
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
	err = service.RemoveOverlay(ctx, entry.Pattern, entry.RelativePath)
	if err != nil {
		t.Fatalf("RemoveOverlay failed: %v", err)
	}

	// Verify it's gone
	_, err = service.GetOverlayContent(ctx, entry)
	if err == nil {
		t.Error("expected error for removed entry, but got nil")
	}

	// Test removing non-existent entry
	err = service.RemoveOverlay(ctx, "non-existent", "file.txt")
	if err == nil {
		t.Error("expected error for non-existent entry, but got nil")
	}
}

// TestApplyOverlays tests applying overlays to a repository
func TestApplyOverlays(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "overlay-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directories for overlays and repository
	repoDir := filepath.Join(tmpDir, "repos", "github.com", "user", "repo")

	// Create repository directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	// Initialize mockWFS for overlay storage
	mockWFS := wfs_mock.NewMockWFS()

	// Create overlay service
	service, err := testtarget.NewOverlayService(mockWFS)
	if err != nil {
		t.Fatalf("NewOverlayService failed: %v", err)
	}

	ctx := context.Background()

	// Test repository
	repo := repository.NewReference("github.com", "user", "repo")

	// Add overlays
	overlays := []struct {
		entry    workspace.OverlayEntry
		content  string
		expected bool // Whether it should be applied to the test repository
	}{
		{
			entry:    workspace.OverlayEntry{Pattern: "github.com/user/repo", RelativePath: ".envrc"},
			content:  "GOPATH=/go",
			expected: true,
		},
		{
			entry:    workspace.OverlayEntry{Pattern: "github.com/user/*", RelativePath: "config/settings.json"},
			content:  `{"debug": true}`,
			expected: true,
		},
		{
			entry:    workspace.OverlayEntry{Pattern: "github.com/other/*", RelativePath: "ignore.txt"},
			content:  "should not be applied",
			expected: false,
		},
	}

	// Add overlays
	for _, o := range overlays {
		err := service.AddOverlay(ctx, o.entry, strings.NewReader(o.content))
		if err != nil {
			t.Fatalf("AddOverlay failed: %v", err)
		}
	}

	// Apply overlays
	err = service.ApplyOverlays(ctx, repo, repoDir)
	if err != nil {
		t.Fatalf("ApplyOverlays failed: %v", err)
	}

	// Verify that expected files were created
	for _, o := range overlays {
		filePath := filepath.Join(repoDir, o.entry.RelativePath)

		if o.expected {
			// Check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("expected file was not created: %s", filePath)
				continue
			}

			// Check if the content is correct
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("failed to read file %s: %v", filePath, err)
				continue
			}

			if string(content) != o.content {
				t.Errorf("content mismatch for %s: got %q, want %q",
					filePath, string(content), o.content)
			}
		} else {
			// Verify that files that shouldn't be applied don't exist
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				t.Errorf("unexpected file was created: %s", filePath)
			}
		}
	}
}
