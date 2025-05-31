package filesystem_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kyoh86/gogh/v4/core/workspace"
	testtarget "github.com/kyoh86/gogh/v4/infra/filesystem"
)

func TestNewOverlayService(t *testing.T) {
	service := testtarget.NewOverlayService()
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if len(service.GetPatterns()) != 0 {
		t.Errorf("Expected empty patterns, got %v", service.GetPatterns())
	}

	if service.HasChanges() {
		t.Error("Expected no changes for new service")
	}
}

func TestAddPattern(t *testing.T) {
	service := testtarget.NewOverlayService()

	// Test adding valid pattern
	files := []workspace.OverlayFile{
		{SourcePath: "/src/file1.go", TargetPath: "dest/file1.go"},
		{SourcePath: "/src/file2.go", TargetPath: "dest/file2.go"},
	}

	if err := service.AddPattern("github.com/*", files); err != nil {
		t.Errorf("Error adding valid pattern: %v", err)
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes after adding pattern")
	}

	patterns := service.GetPatterns()
	if len(patterns) != 1 {
		t.Fatalf("Expected 1 pattern, got %d", len(patterns))
	}

	if patterns[0].Pattern != "github.com/*" {
		t.Errorf("Expected pattern 'github.com/*', got '%s'", patterns[0].Pattern)
	}

	if !reflect.DeepEqual(patterns[0].Files, files) {
		t.Errorf("Expected files %v, got %v", files, patterns[0].Files)
	}

	// Test adding invalid pattern
	if err := service.AddPattern("[invalid", []workspace.OverlayFile{}); err == nil {
		t.Error("Expected error when adding invalid pattern, got nil")
	}

	// Test updating existing pattern
	updatedFiles := []workspace.OverlayFile{
		{SourcePath: "/src/updated.go", TargetPath: "dest/updated.go"},
	}

	service.MarkSaved()
	if service.HasChanges() {
		t.Error("Expected no changes after MarkSaved")
	}

	if err := service.AddPattern("github.com/*", updatedFiles); err != nil {
		t.Errorf("Error updating existing pattern: %v", err)
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes after updating pattern")
	}

	patterns = service.GetPatterns()
	if len(patterns) != 1 {
		t.Fatalf("Expected 1 pattern after update, got %d", len(patterns))
	}

	if !reflect.DeepEqual(patterns[0].Files, updatedFiles) {
		t.Errorf("Expected updated files %v, got %v", updatedFiles, patterns[0].Files)
	}
}

func TestRemovePattern(t *testing.T) {
	service := testtarget.NewOverlayService()

	// Add two patterns
	files1 := []workspace.OverlayFile{{SourcePath: "/src/file1.go", TargetPath: "dest/file1.go"}}
	files2 := []workspace.OverlayFile{{SourcePath: "/src/file2.go", TargetPath: "dest/file2.go"}}

	_ = service.AddPattern("github.com/*", files1)
	_ = service.AddPattern("gitlab.com/*", files2)
	service.MarkSaved()

	// Remove first pattern
	if err := service.RemovePattern("github.com/*"); err != nil {
		t.Errorf("Error removing pattern: %v", err)
	}

	if !service.HasChanges() {
		t.Error("Expected service to have changes after removing pattern")
	}

	patterns := service.GetPatterns()
	if len(patterns) != 1 {
		t.Fatalf("Expected 1 pattern after removal, got %d", len(patterns))
	}

	if patterns[0].Pattern != "gitlab.com/*" {
		t.Errorf("Expected remaining pattern 'gitlab.com/*', got '%s'", patterns[0].Pattern)
	}

	// Try to remove non-existent pattern
	if err := service.RemovePattern("nonexistent"); err == nil {
		t.Error("Expected error when removing non-existent pattern, got nil")
	}
}

func TestGetPatterns(t *testing.T) {
	service := testtarget.NewOverlayService()

	// Initially empty
	if len(service.GetPatterns()) != 0 {
		t.Errorf("Expected empty patterns initially, got %v", service.GetPatterns())
	}

	// Add patterns
	files1 := []workspace.OverlayFile{{SourcePath: "/src/file1.go", TargetPath: "dest/file1.go"}}
	files2 := []workspace.OverlayFile{{SourcePath: "/src/file2.go", TargetPath: "dest/file2.go"}}

	_ = service.AddPattern("github.com/*", files1)
	_ = service.AddPattern("gitlab.com/*", files2)

	patterns := service.GetPatterns()
	if len(patterns) != 2 {
		t.Fatalf("Expected 2 patterns, got %d", len(patterns))
	}

	// Check that GetPatterns returns a copy (not the original slice)
	// by modifying the returned slice and checking that the service is unaffected
	patterns[0].Pattern = "modified"

	originalPatterns := service.GetPatterns()
	if originalPatterns[0].Pattern == "modified" {
		t.Error("GetPatterns should return a copy, not the original slice")
	}
}

func TestApplyToRepository(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "overlay-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source files
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create two source files with content
	file1Path := filepath.Join(srcDir, "file1.txt")
	if err := os.WriteFile(file1Path, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create source file1: %v", err)
	}

	file2Path := filepath.Join(srcDir, "file2.txt")
	if err := os.WriteFile(file2Path, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create source file2: %v", err)
	}

	// Create repository directory
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Set up the service
	service := testtarget.NewOverlayService()

	// Add pattern matching our test repository
	files := []workspace.OverlayFile{
		{SourcePath: file1Path, TargetPath: "subdir/dest1.txt"},
		{SourcePath: file2Path, TargetPath: "dest2.txt"},
	}

	_ = service.AddPattern("github.com/kyoh86/gogh", files)

	// Apply overlays to repository
	ctx := context.Background()
	if err := service.ApplyToRepository(ctx, repoDir, "github.com/kyoh86/gogh"); err != nil {
		t.Fatalf("ApplyToRepository failed: %v", err)
	}

	// Verify files were copied to the right locations
	dest1Path := filepath.Join(repoDir, "subdir", "dest1.txt")
	dest1Content, err := os.ReadFile(dest1Path)
	if err != nil {
		t.Fatalf("Failed to read destination file1: %v", err)
	}
	if string(dest1Content) != "content1" {
		t.Errorf("Expected content1, got %s", string(dest1Content))
	}

	dest2Path := filepath.Join(repoDir, "dest2.txt")
	dest2Content, err := os.ReadFile(dest2Path)
	if err != nil {
		t.Fatalf("Failed to read destination file2: %v", err)
	}
	if string(dest2Content) != "content2" {
		t.Errorf("Expected content2, got %s", string(dest2Content))
	}

	// Test skip existing files
	// First modify the existing file
	if err := os.WriteFile(dest2Path, []byte("modified"), 0644); err != nil {
		t.Fatalf("Failed to modify destination file: %v", err)
	}

	// Apply again - should skip the existing file
	if err := service.ApplyToRepository(ctx, repoDir, "github.com/kyoh86/gogh"); err != nil {
		t.Fatalf("Second ApplyToRepository failed: %v", err)
	}

	// Verify the file wasn't overwritten
	dest2ContentAfter, err := os.ReadFile(dest2Path)
	if err != nil {
		t.Fatalf("Failed to read destination file2 after second apply: %v", err)
	}
	if string(dest2ContentAfter) != "modified" {
		t.Errorf("Expected content to remain 'modified', got %s", string(dest2ContentAfter))
	}
}

func TestApplyToRepository_NonMatchingRepo(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "overlay-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up the service
	service := testtarget.NewOverlayService()

	// Add pattern that won't match our test repository
	files := []workspace.OverlayFile{
		{SourcePath: "/non/existent/file.txt", TargetPath: "dest.txt"},
	}

	_ = service.AddPattern("github.com/other/*", files)

	// Apply overlays to repository (should do nothing since no patterns match)
	ctx := context.Background()
	if err := service.ApplyToRepository(ctx, tmpDir, "github.com/kyoh86/gogh"); err != nil {
		t.Fatalf("ApplyToRepository failed: %v", err)
	}

	// Verify no files were created
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}
	if len(entries) > 0 {
		t.Errorf("Expected no files to be created, got %d entries", len(entries))
	}
}

func TestApplyToRepository_ErrorCases(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "overlay-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Set up the service
	service := testtarget.NewOverlayService()

	// Add pattern with non-existent source file
	files := []workspace.OverlayFile{
		{SourcePath: "/non/existent/file.txt", TargetPath: "dest.txt"},
	}

	_ = service.AddPattern("github.com/kyoh86/gogh", files)

	// Apply overlays to repository (should fail because source file doesn't exist)
	ctx := context.Background()
	if err := service.ApplyToRepository(ctx, tmpDir, "github.com/kyoh86/gogh"); err == nil {
		t.Error("Expected error when source file doesn't exist, got nil")
	}
}

func TestMarkSavedAndHasChanges(t *testing.T) {
	service := testtarget.NewOverlayService()

	// Initially has no changes
	if service.HasChanges() {
		t.Error("New service should not have changes")
	}

	// Add a pattern to create changes
	files := []workspace.OverlayFile{{SourcePath: "/src/file.go", TargetPath: "dest/file.go"}}
	_ = service.AddPattern("github.com/*", files)

	if !service.HasChanges() {
		t.Error("Expected service to have changes after adding pattern")
	}

	// Mark as saved
	service.MarkSaved()

	if service.HasChanges() {
		t.Error("Expected no changes after MarkSaved")
	}

	// Remove pattern to create changes again
	_ = service.RemovePattern("github.com/*")

	if !service.HasChanges() {
		t.Error("Expected service to have changes after removing pattern")
	}
}
