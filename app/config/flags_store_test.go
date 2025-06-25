package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
)

func TestFlagsStore_Source(t *testing.T) {
	// Save and restore environment variable
	oldFlagPath := os.Getenv("GOGH_FLAG_PATH")
	defer os.Setenv("GOGH_FLAG_PATH", oldFlagPath)

	t.Run("default path", func(t *testing.T) {
		os.Unsetenv("GOGH_FLAG_PATH")
		store := config.NewFlagsStore()
		path, err := store.Source()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filepath.Base(path) != "flags.v4.toml" {
			t.Errorf("expected path to contain flags.v4.toml, got %s", path)
		}
	})

	t.Run("custom path from env", func(t *testing.T) {
		customPath := "/custom/path/flags.toml"
		os.Setenv("GOGH_FLAG_PATH", customPath)
		store := config.NewFlagsStore()
		path, err := store.Source()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != customPath {
			t.Errorf("expected path %s, got %s", customPath, path)
		}
	})
}

func TestFlagsStore_Load(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "flags_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test flags file
	flagsContent := `
[list]
limit = 200
format = "json"
primary = true

[repos]
limit = 50
privacy = "public"
fork = "exclude"
`
	flagsPath := filepath.Join(tempDir, "flags.v4.toml")
	err = os.WriteFile(flagsPath, []byte(flagsContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Save and restore environment variable
	oldFlagPath := os.Getenv("GOGH_FLAG_PATH")
	defer os.Setenv("GOGH_FLAG_PATH", oldFlagPath)

	t.Run("successful load", func(t *testing.T) {
		os.Setenv("GOGH_FLAG_PATH", flagsPath)
		store := config.NewFlagsStore()
		initial := config.DefaultFlags

		flags, err := store.Load(context.Background(), initial)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check List flags
		if flags.List.Limit != 200 {
			t.Errorf("expected List.Limit to be 200, got %d", flags.List.Limit)
		}
		if flags.List.Format != "json" {
			t.Errorf("expected List.Format to be 'json', got '%s'", flags.List.Format)
		}
		if !flags.List.Primary {
			t.Errorf("expected List.Primary to be true")
		}

		// Check Repos flags
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
		nonExistentPath := filepath.Join(tempDir, "nonexistent.toml")
		os.Setenv("GOGH_FLAG_PATH", nonExistentPath)

		store := config.NewFlagsStore()
		initial := config.DefaultFlags

		_, err := store.Load(context.Background(), initial)
		if err == nil {
			t.Error("expected error when file not found, got nil")
		}
	})

	t.Run("invalid toml", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid.toml")
		err = os.WriteFile(invalidPath, []byte("invalid toml content"), 0o644)
		if err != nil {
			t.Fatalf("failed to write invalid test file: %v", err)
		}

		os.Setenv("GOGH_FLAG_PATH", invalidPath)
		store := config.NewFlagsStore()
		initial := config.DefaultFlags

		_, err := store.Load(context.Background(), initial)
		if err == nil {
			t.Error("expected error with invalid TOML, got nil")
		}
	})
}

func TestFlagsStore_Save(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "flags_save_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save and restore environment variable
	oldFlagPath := os.Getenv("GOGH_FLAG_PATH")
	defer os.Setenv("GOGH_FLAG_PATH", oldFlagPath)

	// Set up the flags path in the temp directory
	flagsPath := filepath.Join(tempDir, "flags.v4.toml")
	os.Setenv("GOGH_FLAG_PATH", flagsPath)

	// Create store and initial flags
	store := config.NewFlagsStore()
	ctx := context.Background()

	t.Run("save with changes", func(t *testing.T) {
		// Create flags with changes
		flags := config.DefaultFlags()
		flags.List.Limit = 300
		flags.List.Format = "custom"
		flags.Repos.Privacy = "private"

		// Mark as changed
		flags.RawHasChanges = true

		// Save the flags
		err := store.Save(ctx, flags, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(flagsPath); os.IsNotExist(err) {
			t.Fatal("flags file was not created")
		}

		// Load the saved flags and verify content
		loadedFlags, err := store.Load(ctx, config.DefaultFlags)
		if err != nil {
			t.Fatalf("failed to load saved flags: %v", err)
		}

		// Check values were correctly saved
		if loadedFlags.List.Limit != 300 {
			t.Errorf("expected List.Limit to be 300, got %d", loadedFlags.List.Limit)
		}
		if loadedFlags.List.Format != "custom" {
			t.Errorf("expected List.Format to be 'custom', got '%s'", loadedFlags.List.Format)
		}
		if loadedFlags.Repos.Privacy != "private" {
			t.Errorf("expected Repos.Privacy to be 'private', got '%s'", loadedFlags.Repos.Privacy)
		}
	})

	t.Run("no save when no changes", func(t *testing.T) {
		// Delete the existing file
		os.Remove(flagsPath)

		// Create flags with no changes
		flags := config.DefaultFlags()
		flags.MarkSaved() // Explicitly mark as saved (no changes)

		// Attempt to save
		err := store.Save(ctx, flags, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file was not created
		if _, err := os.Stat(flagsPath); !os.IsNotExist(err) {
			t.Fatal("flags file was created when it shouldn't have been")
		}
	})

	t.Run("save with force", func(t *testing.T) {
		// Delete the existing file
		os.Remove(flagsPath)

		// Create flags with no changes
		flags := config.DefaultFlags()
		flags.MarkSaved() // Explicitly mark as saved (no changes)

		// Save with force=true
		err := store.Save(ctx, flags, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file was created despite no changes
		if _, err := os.Stat(flagsPath); os.IsNotExist(err) {
			t.Fatal("flags file was not created when forced")
		}
	})

	t.Run("does not error on write to non-exist directory", func(t *testing.T) {
		// Set an invalid path (not a directory)
		invalidPath := filepath.Join(tempDir, "non-exist", "dir", "flags.toml")
		os.Setenv("GOGH_FLAG_PATH", invalidPath)

		// Create flags with changes
		flags := config.DefaultFlags()
		flags.RawHasChanges = true

		// Try to save to an invalid directory
		store := config.NewFlagsStore() // Need a new store to pick up the env var change
		err := store.Save(ctx, flags, true)
		// Should not fail to create directory
		if err != nil {
			t.Fatalf("expected success when saving to invalid directory, got %v", err)
		}
	})
}
