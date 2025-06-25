package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/infra/filesystem"
)

func TestWorkspaceService_GetRoots(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Test empty roots
	roots := ws.GetRoots()
	if len(roots) != 0 {
		t.Errorf("expected 0 roots, got %d", len(roots))
	}

	// Add roots
	root1, _ := filepath.Abs("/home/user/repos1")
	root2, _ := filepath.Abs("/home/user/repos2")

	err := ws.AddRoot(root1, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	err = ws.AddRoot(root2, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Test roots after adding
	roots = ws.GetRoots()
	if len(roots) != 2 {
		t.Errorf("expected 2 roots, got %d", len(roots))
	}
}

func TestWorkspaceService_GetPrimaryRoot(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Test no primary root
	primaryRoot := ws.GetPrimaryRoot()
	if primaryRoot != "" {
		t.Errorf("expected empty primary root, got %q", primaryRoot)
	}

	// Add first root (should become primary)
	root1, _ := filepath.Abs("/home/user/repos1")
	root2, _ := filepath.Abs("/home/user/repos2")

	err := ws.AddRoot(root1, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	primaryRoot = ws.GetPrimaryRoot()
	if primaryRoot != root1 {
		t.Errorf("expected primary root to be %s, got %q", root1, primaryRoot)
	}

	// Add second root as primary
	err = ws.AddRoot(root2, true)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	primaryRoot = ws.GetPrimaryRoot()
	if primaryRoot != root2 {
		t.Errorf("expected primary root to be %s, got %q", root2, primaryRoot)
	}
}

func TestWorkspaceService_SetPrimaryRoot(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add roots
	root1, _ := filepath.Abs("/home/user/repos1")
	root2, _ := filepath.Abs("/home/user/repos2")

	err := ws.AddRoot(root1, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	err = ws.AddRoot(root2, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Set primary root
	err = ws.SetPrimaryRoot(root2)
	if err != nil {
		t.Fatalf("failed to set primary root: %v", err)
	}

	primaryRoot := ws.GetPrimaryRoot()
	if primaryRoot != root2 {
		t.Errorf("expected primary root to be %s, got %q", root2, primaryRoot)
	}

	// Try to set non-existent root as primary
	root3, _ := filepath.Abs("/home/user/repos3")
	err = ws.SetPrimaryRoot(root3)
	if err == nil {
		t.Error("expected error when setting non-existent root as primary")
	}
}

func TestWorkspaceService_AddRoot(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add root
	root1, _ := filepath.Abs("/home/user/repos1")

	err := ws.AddRoot(root1, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Try to add duplicate root
	err = ws.AddRoot(root1, false)
	if err == nil {
		t.Error("expected error when adding duplicate root")
	}

	// Add relative path (should be converted to absolute)
	err = ws.AddRoot("./repos2", false)
	if err != nil {
		t.Fatalf("failed to add relative root: %v", err)
	}

	roots := ws.GetRoots()
	if len(roots) != 2 {
		t.Errorf("expected 2 roots, got %d", len(roots))
	}
}

func TestWorkspaceService_RemoveRoot(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add roots
	root1, _ := filepath.Abs("/home/user/repos1")
	root2, _ := filepath.Abs("/home/user/repos2")

	err := ws.AddRoot(root1, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	err = ws.AddRoot(root2, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Remove root
	err = ws.RemoveRoot(root2)
	if err != nil {
		t.Fatalf("failed to remove root: %v", err)
	}

	roots := ws.GetRoots()
	if len(roots) != 1 {
		t.Errorf("expected 1 root, got %d", len(roots))
	}

	// Remove non-existent root
	root3, _ := filepath.Abs("/home/user/repos3")
	err = ws.RemoveRoot(root3)
	if err == nil {
		t.Error("expected error when removing non-existent root")
	}

	// Remove primary root
	err = ws.RemoveRoot(root1)
	if err != nil {
		t.Fatalf("failed to remove primary root: %v", err)
	}

	primaryRoot := ws.GetPrimaryRoot()
	if primaryRoot != "" {
		t.Errorf("expected empty primary root after removing all roots, got %q", primaryRoot)
	}
}

func TestWorkspaceService_RemovePrimaryRoot(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add multiple roots
	root1, _ := filepath.Abs("/home/user/repos1")
	root2, _ := filepath.Abs("/home/user/repos2")
	root3, _ := filepath.Abs("/home/user/repos3")

	err := ws.AddRoot(root1, true)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	err = ws.AddRoot(root2, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	err = ws.AddRoot(root3, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Verify primary root
	primaryRoot := ws.GetPrimaryRoot()
	if primaryRoot != root1 {
		t.Errorf("expected primary root to be %s, got %q", root1, primaryRoot)
	}

	// Remove primary root
	err = ws.RemoveRoot(root1)
	if err != nil {
		t.Fatalf("failed to remove primary root: %v", err)
	}

	// Primary root should switch to the first remaining root
	primaryRoot = ws.GetPrimaryRoot()
	if primaryRoot != root2 {
		t.Errorf("expected primary root to be %s after removal, got %q", root2, primaryRoot)
	}
}

func TestWorkspaceService_GetLayoutFor(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add root
	root, _ := filepath.Abs("/home/user/repos")

	err := ws.AddRoot(root, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Get layout for root
	layout := ws.GetLayoutFor(root)
	if layout == nil {
		t.Fatal("expected layout, got nil")
	}

	if layout.GetRoot() != root {
		t.Errorf("expected layout root to be %s, got %q", root, layout.GetRoot())
	}
}

func TestWorkspaceService_GetPrimaryLayout(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Add root
	root, _ := filepath.Abs("/home/user/repos")

	err := ws.AddRoot(root, true)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Get primary layout
	layout := ws.GetPrimaryLayout()
	if layout == nil {
		t.Fatal("expected layout, got nil")
	}

	if layout.GetRoot() != root {
		t.Errorf("expected layout root to be %s, got %q", root, layout.GetRoot())
	}
}

func TestWorkspaceService_HasChanges(t *testing.T) {
	ws := filesystem.NewWorkspaceService()

	// Initially no changes
	if ws.HasChanges() {
		t.Error("expected no changes initially")
	}

	// Add root
	root, _ := filepath.Abs("/home/user/repos")

	err := ws.AddRoot(root, false)
	if err != nil {
		t.Fatalf("failed to add root: %v", err)
	}

	// Should have changes after adding root
	if !ws.HasChanges() {
		t.Error("expected changes after adding root")
	}

	// Mark saved
	ws.MarkSaved()
	if ws.HasChanges() {
		t.Error("expected no changes after marking saved")
	}

	// Set primary root
	err = ws.SetPrimaryRoot(root)
	if err != nil {
		t.Fatalf("failed to set primary root: %v", err)
	}

	// Should have changes after setting primary root
	if !ws.HasChanges() {
		t.Error("expected changes after setting primary root")
	}
}

func TestWorkspaceService_Errors(t *testing.T) {
	t.Run("ErrRootNotFound", func(t *testing.T) {
		err := workspace.ErrRootNotFound
		if err.Error() != "repository root not found" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("ErrRootAlreadyExists", func(t *testing.T) {
		err := workspace.ErrRootAlreadyExists
		if err.Error() != "repository root already exists" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("ErrNoPrimaryRoot", func(t *testing.T) {
		err := workspace.ErrNoPrimaryRoot
		if err.Error() != "no primary repository root configured" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
