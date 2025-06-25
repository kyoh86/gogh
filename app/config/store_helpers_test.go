package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTOMLFile(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		// Create temporary file with valid TOML
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.toml")

		type TestData struct {
			Name  string `toml:"name"`
			Value int    `toml:"value"`
		}

		content := `name = "test"
value = 42`
		if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Test loading
		data, err := loadTOMLFile[TestData](testFile)
		if err != nil {
			t.Fatalf("loadTOMLFile failed: %v", err)
		}

		if data.Name != "test" || data.Value != 42 {
			t.Errorf("Unexpected data: got %+v", data)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		type TestData struct {
			Name string `toml:"name"`
		}

		_, err := loadTOMLFile[TestData]("/nonexistent/file.toml")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
		if !os.IsNotExist(err) {
			t.Errorf("Expected os.IsNotExist error, got: %v", err)
		}
	})

	t.Run("invalid TOML", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "invalid.toml")

		// Write invalid TOML
		if err := os.WriteFile(testFile, []byte("invalid [toml content"), 0o644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		type TestData struct {
			Name string `toml:"name"`
		}

		_, err := loadTOMLFile[TestData](testFile)
		if err == nil {
			t.Error("Expected error for invalid TOML")
		}
	})
}

func TestSaveTOMLFile(t *testing.T) {
	t.Run("successful save", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "output.toml")

		type TestData struct {
			Name  string `toml:"name"`
			Value int    `toml:"value"`
		}

		data := TestData{
			Name:  "test",
			Value: 42,
		}

		err := saveTOMLFile(testFile, data)
		if err != nil {
			t.Fatalf("saveTOMLFile failed: %v", err)
		}

		// Verify file exists and content is correct
		content, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		expectedContent := `name = 'test'
value = 42
`
		if string(content) != expectedContent {
			t.Errorf("Unexpected content:\ngot: %s\nwant: %s", content, expectedContent)
		}
	})

	t.Run("creates directory if not exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "subdir", "nested", "output.toml")

		type TestData struct {
			Name string `toml:"name"`
		}

		data := TestData{Name: "test"}

		err := saveTOMLFile(testFile, data)
		if err != nil {
			t.Fatalf("saveTOMLFile failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(testFile); err != nil {
			t.Errorf("File was not created: %v", err)
		}
	})

	t.Run("permission error", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Cannot test permission errors as root")
		}

		tempDir := t.TempDir()
		// Create a directory with no write permission
		readOnlyDir := filepath.Join(tempDir, "readonly")
		if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
			t.Fatalf("Failed to create readonly dir: %v", err)
		}

		testFile := filepath.Join(readOnlyDir, "output.toml")

		type TestData struct {
			Name string `toml:"name"`
		}

		data := TestData{Name: "test"}

		err := saveTOMLFile(testFile, data)
		if err == nil {
			t.Error("Expected error for permission denied")
		}
	})
}

func TestEnsureDirectoryExists(t *testing.T) {
	t.Run("creates nested directories", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "a", "b", "c", "file.txt")

		err := ensureDirectoryExists(filePath)
		if err != nil {
			t.Fatalf("ensureDirectoryExists failed: %v", err)
		}

		// Verify directory exists
		dirPath := filepath.Dir(filePath)
		if _, err := os.Stat(dirPath); err != nil {
			t.Errorf("Directory was not created: %v", err)
		}
	})

	t.Run("existing directory", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "file.txt")

		// Directory already exists
		err := ensureDirectoryExists(filePath)
		if err != nil {
			t.Errorf("ensureDirectoryExists failed for existing directory: %v", err)
		}
	})

	t.Run("permission error", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Cannot test permission errors as root")
		}

		tempDir := t.TempDir()
		// Create a directory with no write permission
		readOnlyDir := filepath.Join(tempDir, "readonly")
		if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
			t.Fatalf("Failed to create readonly dir: %v", err)
		}

		filePath := filepath.Join(readOnlyDir, "subdir", "file.txt")

		err := ensureDirectoryExists(filePath)
		if err == nil {
			t.Error("Expected error for permission denied")
		}
	})
}

func TestOpenFileForWrite(t *testing.T) {
	t.Run("creates new file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "new.txt")

		file, err := openFileForWrite(testFile)
		if err != nil {
			t.Fatalf("openFileForWrite failed: %v", err)
		}
		defer file.Close()

		// Write some data
		if _, err := file.WriteString("test"); err != nil {
			t.Errorf("Failed to write to file: %v", err)
		}
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "existing.txt")

		// Create existing file with content
		if err := os.WriteFile(testFile, []byte("old content"), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		file, err := openFileForWrite(testFile)
		if err != nil {
			t.Fatalf("openFileForWrite failed: %v", err)
		}
		defer file.Close()

		// Write new content
		if _, err := file.WriteString("new"); err != nil {
			t.Errorf("Failed to write to file: %v", err)
		}

		// Verify file was truncated
		file.Close()
		content, _ := os.ReadFile(testFile)
		if string(content) != "new" {
			t.Errorf("File was not truncated: got %q", content)
		}
	})

	t.Run("permission error", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Cannot test permission errors as root")
		}

		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "readonly.txt")

		// Create file with read-only permission
		if err := os.WriteFile(testFile, []byte("test"), 0o444); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err := openFileForWrite(testFile)
		if err == nil {
			t.Error("Expected error for read-only file")
		}
	})
}
