package wfs_mock

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestMockWFS_DirectoryHierarchy(t *testing.T) {
	mockFS := NewMockWFS()

	// Create a deep directory structure
	err := mockFS.MkdirAll("a/b/c", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory hierarchy: %v", err)
	}

	// Verify all directories exist
	dirs := []string{"a", "a/b", "a/b/c"}
	for _, dir := range dirs {
		stat, err := mockFS.Stat(dir)
		if err != nil {
			t.Errorf("Failed to stat directory %s: %v", dir, err)
			continue
		}
		if !stat.IsDir() {
			t.Errorf("Expected %s to be a directory", dir)
		}
	}

	// Check parent-child relationships
	checkDirContains := func(parent, child string) {
		entries, err := mockFS.ReadDir(parent)
		if err != nil {
			t.Fatalf("Failed to read directory %s: %v", parent, err)
		}

		childName := filepath.Base(child)
		found := false
		for _, entry := range entries {
			if entry.Name() == childName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Directory %s does not contain %s", parent, childName)
		}
	}

	checkDirContains("", "a")
	checkDirContains("a", "a/b")
	checkDirContains("a/b", "a/b/c")
}

func TestMockWFS_FileOperations(t *testing.T) {
	mockFS := NewMockWFS()

	// Create a deep file
	content := []byte("test content")
	err := mockFS.WriteFile("a/b/c/test.txt", content, 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Verify file exists
	stat, err := mockFS.Stat("a/b/c/test.txt")
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.IsDir() {
		t.Error("Expected file but got directory")
	}
	if stat.Size() != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), stat.Size())
	}

	// Read file
	readContent, err := mockFS.ReadFile("a/b/c/test.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if !bytes.Equal(readContent, content) {
		t.Errorf("File content mismatch: expected %s, got %s", content, readContent)
	}

	// Verify directory entries
	checkDirContains := func(dir string, expectedItem string) {
		entries, err := mockFS.ReadDir(dir)
		if err != nil {
			t.Fatalf("Failed to read directory %s: %v", dir, err)
		}

		baseName := filepath.Base(expectedItem)
		found := false
		for _, entry := range entries {
			if entry.Name() == baseName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Directory %s does not contain %s", dir, baseName)
		}
	}

	// Check that each directory contains the expected child
	checkDirContains("", "a")
	checkDirContains("a", "a/b")
	checkDirContains("a/b", "a/b/c")
	checkDirContains("a/b/c", "a/b/c/test.txt")
}

func TestMockWFS_Create(t *testing.T) {
	mockFS := NewMockWFS()

	// Create file using Create method
	writer, err := mockFS.Create("a/b/c/create.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	content := []byte("created file content")
	_, err = writer.Write(content)
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	// Verify file exists and directories were created
	checkDirContains := func(dir string, expectedItem string) {
		entries, err := mockFS.ReadDir(dir)
		if err != nil {
			t.Fatalf("Failed to read directory %s: %v", dir, err)
		}

		baseName := filepath.Base(expectedItem)
		found := false
		for _, entry := range entries {
			if entry.Name() == baseName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Directory %s does not contain %s", dir, baseName)
		}
	}

	// Read back and verify content
	readContent, err := mockFS.ReadFile("a/b/c/create.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if !bytes.Equal(readContent, content) {
		t.Errorf("File content mismatch: expected %s, got %s", content, readContent)
	}

	// Check directory hierarchy
	checkDirContains("", "a")
	checkDirContains("a", "a/b")
	checkDirContains("a/b", "a/b/c")
	checkDirContains("a/b/c", "a/b/c/create.txt")
}
