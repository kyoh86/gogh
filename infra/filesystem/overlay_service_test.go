package filesystem_test

import (
	"context"
	"io/fs"
	"strings"
	"testing"

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
