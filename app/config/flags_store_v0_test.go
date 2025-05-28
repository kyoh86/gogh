package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
)

func TestFlagsStoreV0_Source(t *testing.T) {
	// Save and restore environment variable
	oldFlagPath := os.Getenv("GOGH_FLAG_PATH")
	defer os.Setenv("GOGH_FLAG_PATH", oldFlagPath)

	t.Run("default path", func(t *testing.T) {
		os.Unsetenv("GOGH_FLAG_PATH")
		store := testtarget.NewFlagsStoreV0()
		path, err := store.Source()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filepath.Base(path) != "flag.yaml" {
			t.Errorf("expected path to contain flag.yaml, got %s", path)
		}
	})

	t.Run("custom path from env", func(t *testing.T) {
		customPath := "/custom/path/flag.yaml"
		os.Setenv("GOGH_FLAG_PATH", customPath)
		store := testtarget.NewFlagsStoreV0()
		path, err := store.Source()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != customPath {
			t.Errorf("expected path %s, got %s", customPath, path)
		}
	})
}

func TestFlagsStoreV0_Load(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "flags_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test flags file
	flagsContent := `
list:
  limit: 200
  format: "json"
  primary: true
repos:
  limit: 50
  privacy: "public"
  fork: "exclude"
`
	flagsPath := filepath.Join(tempDir, "flag.yaml")
	err = os.WriteFile(flagsPath, []byte(flagsContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Save and restore environment variable
	oldFlagPath := os.Getenv("GOGH_FLAG_PATH")
	defer os.Setenv("GOGH_FLAG_PATH", oldFlagPath)

	t.Run("successful load", func(t *testing.T) {
		os.Setenv("GOGH_FLAG_PATH", flagsPath)
		store := testtarget.NewFlagsStoreV0()
		initial := testtarget.DefaultFlags

		flags, err := store.Load(context.Background(), initial)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if flags.List.Limit != 200 {
			t.Errorf("expected List.Limit to be 200, got %d", flags.List.Limit)
		}
		if flags.List.Format != "json" {
			t.Errorf("expected List.Format to be 'json', got '%s'", flags.List.Format)
		}
		if !flags.List.Primary {
			t.Errorf("expected List.Primary to be true")
		}

		if flags.Repos.Limit != 50 {
			t.Errorf("expected Repos.Limit to be 50, got %d", flags.Repos.Limit)
		}
		if flags.Repos.Privacy != "public" {
			t.Errorf("expected Repos.Privacy to be 'public', got '%s'", flags.Repos.Privacy)
		}
		if flags.Repos.Fork != "exclude" {
			t.Errorf("expected Repos.Fork to be 'exclude', got '%s'", flags.Repos.Fork)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.yaml")
		os.Setenv("GOGH_FLAG_PATH", nonExistentPath)

		store := testtarget.NewFlagsStoreV0()
		initial := testtarget.DefaultFlags

		_, err := store.Load(context.Background(), initial)
		if err == nil {
			t.Error("expected error when file not found, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid.yaml")
		err = os.WriteFile(invalidPath, []byte("invalid: yaml: content: - :\n"), 0644)
		if err != nil {
			t.Fatalf("failed to write invalid test file: %v", err)
		}

		os.Setenv("GOGH_FLAG_PATH", invalidPath)
		store := testtarget.NewFlagsStoreV0()
		initial := testtarget.DefaultFlags

		_, err := store.Load(context.Background(), initial)
		if err == nil {
			t.Error("expected error with invalid YAML, got nil")
		}
	})
}
