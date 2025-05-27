package filesystem_test

import (
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/core/workspace"
	testtarget "github.com/kyoh86/gogh/v4/infra/filesystem"
)

func TestNewWorkspaceService(t *testing.T) {
	service := testtarget.NewWorkspaceService()
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if len(service.GetRoots()) != 0 {
		t.Errorf("Expected empty roots, got %v", service.GetRoots())
	}

	if service.GetPrimaryRoot() != "" {
		t.Errorf("Expected empty primary root, got %v", service.GetPrimaryRoot())
	}

	if service.HasChanges() {
		t.Error("Expected no changes for new service")
	}
}

func TestAddRoot(t *testing.T) {
	service := testtarget.NewWorkspaceService()

	root1, err := filepath.Abs("/test/path1")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	root2, err := filepath.Abs("/test/path2")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	root3, err := filepath.Abs("/test/path3")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}

	// Test adding first root (should become primary)
	if err := service.AddRoot(root1, false); err != nil {
		t.Errorf("Error adding root: %v", err)
	}

	if len(service.GetRoots()) != 1 {
		t.Errorf("Expected 1 root, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != root1 {
		t.Errorf("Expected primary root %s, got %s", root1, service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Reset changed flag
	service.MarkSaved()
	if service.HasChanges() {
		t.Error("Expected no changes after MarkSaved")
	}

	// Test adding second root (should not change primary)
	if err := service.AddRoot(root2, false); err != nil {
		t.Errorf("Error adding root: %v", err)
	}

	if len(service.GetRoots()) != 2 {
		t.Errorf("Expected 2 roots, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != root1 {
		t.Errorf("Expected primary root to remain %s, got %s", root1, service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Test adding third root as primary
	if err := service.AddRoot(root3, true); err != nil {
		t.Errorf("Error adding root: %v", err)
	}

	if len(service.GetRoots()) != 3 {
		t.Errorf("Expected 3 roots, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != root3 {
		t.Errorf("Expected primary root %s, got %s", root3, service.GetPrimaryRoot())
	}

	// Test adding duplicate root
	if err := service.AddRoot(root1, false); err == nil {
		t.Error("Expected error when adding duplicate root, got nil")
	}

	if len(service.GetRoots()) != 3 {
		t.Errorf("Expected still 3 roots, got %d", len(service.GetRoots()))
	}
}

func TestSetPrimaryRoot(t *testing.T) {
	service := testtarget.NewWorkspaceService()

	root1, err := filepath.Abs("/test/path1")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	root2, err := filepath.Abs("/test/path2")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}

	// Add two roots
	_ = service.AddRoot(root1, false)
	_ = service.AddRoot(root2, false)
	service.MarkSaved()

	// Set second root as primary
	if err := service.SetPrimaryRoot(root2); err != nil {
		t.Errorf("Error setting primary root: %v", err)
	}

	if service.GetPrimaryRoot() != root2 {
		t.Errorf("Expected primary root %s, got %s", root2, service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Try to set non-existent root as primary
	if err := service.SetPrimaryRoot("/test/nonexistent"); err == nil {
		t.Error("Expected error when setting non-existent root as primary, got nil")
	}

	if service.GetPrimaryRoot() != root2 {
		t.Errorf("Expected primary root to remain %s, got %s", root2, service.GetPrimaryRoot())
	}
}

func TestRemoveRoot(t *testing.T) {
	service := testtarget.NewWorkspaceService()

	root1, err := filepath.Abs("/test/path1")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	root2, err := filepath.Abs("/test/path2")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	root3, err := filepath.Abs("/test/path3")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}

	// Add three roots
	_ = service.AddRoot(root1, false)
	_ = service.AddRoot(root2, false)
	_ = service.AddRoot(root3, false)

	// Set second root as primary
	_ = service.SetPrimaryRoot(root2)
	service.MarkSaved()

	// Remove a non-primary root
	if err := service.RemoveRoot(root3); err != nil {
		t.Errorf("Error removing root: %v", err)
	}

	if len(service.GetRoots()) != 2 {
		t.Errorf("Expected 2 roots, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != root2 {
		t.Errorf("Expected primary root to remain %s, got %s", root2, service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Remove primary root
	service.MarkSaved()
	if err := service.RemoveRoot(root2); err != nil {
		t.Errorf("Error removing primary root: %v", err)
	}

	if len(service.GetRoots()) != 1 {
		t.Errorf("Expected 1 root, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != root1 {
		t.Errorf("Expected new primary root %s, got %s", root1, service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Remove last root
	service.MarkSaved()
	if err := service.RemoveRoot(root1); err != nil {
		t.Errorf("Error removing last root: %v", err)
	}

	if len(service.GetRoots()) != 0 {
		t.Errorf("Expected 0 roots, got %d", len(service.GetRoots()))
	}

	if service.GetPrimaryRoot() != "" {
		t.Errorf("Expected empty primary root, got %s", service.GetPrimaryRoot())
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes")
	}

	// Try to remove non-existent root
	service.MarkSaved()
	if err := service.RemoveRoot("/test/nonexistent"); err == nil {
		t.Error("Expected error when removing non-existent root, got nil")
	}

	if service.HasChanges() {
		t.Error("Expected no changes after failed removal")
	}
}

func TestGetLayoutFor(t *testing.T) {
	service := testtarget.NewWorkspaceService()
	root1, err := filepath.Abs("/test/path1")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	_ = service.AddRoot(root1, false)

	layout := service.GetLayoutFor(root1)
	if layout == nil {
		t.Error("Expected non-nil layout")
	}
}

func TestGetPrimaryLayout(t *testing.T) {
	service := testtarget.NewWorkspaceService()
	root1, err := filepath.Abs("/test/path1")
	if err != nil {
		t.Fatalf("Error getting absolute path: %v", err)
	}
	_ = service.AddRoot(root1, false)

	layout := service.GetPrimaryLayout()
	if layout == nil {
		t.Error("Expected non-nil primary layout")
	}
}
