package config_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
)

func TestOverlayContentStore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "overlay-content-store-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock AppContextPathFunc to use our temp directory
	origAppContextPathFunc := config.AppContextPathFunc
	defer func() { config.AppContextPathFunc = origAppContextPathFunc }()

	config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return filepath.Join(tempDir, "overlay.v4"), nil
	}

	// Create test context and store
	ctx := context.Background()
	store := config.NewOverlayContentStore()

	t.Run("Source", func(t *testing.T) {
		source, err := store.Source()
		if err != nil {
			t.Fatalf("Source() failed: %v", err)
		}
		expected := filepath.Join(tempDir, "overlay.v4")
		if source != expected {
			t.Errorf("Source() = %q, want %q", source, expected)
		}
	})

	// Create test data
	testContent := []byte("test content data")
	const testID = "test-overlay-id"
	t.Run("Save", func(t *testing.T) {
		buffer := bytes.NewBuffer(testContent)
		if err := store.Save(ctx, testID, buffer); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		// Verify the content was saved to the expected path
		source, _ := store.Source()
		savedPath := filepath.Join(source, testID)
		if _, err := os.Stat(savedPath); os.IsNotExist(err) {
			t.Errorf("SaveContent did not create file at %q", savedPath)
		}

		// Verify the content is correct
		savedContent, err := os.ReadFile(savedPath)
		if err != nil {
			t.Fatalf("Failed to read saved content: %v", err)
		}
		if !bytes.Equal(savedContent, testContent) {
			t.Errorf("Saved content doesn't match original content")
		}
	})

	t.Run("Open", func(t *testing.T) {
		reader, err := store.Open(ctx, testID)
		if err != nil {
			t.Fatalf("Open() failed: %v", err)
		}
		defer reader.Close()

		// Read the content
		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("Failed to read content: %v", err)
		}

		// Verify content matches what we saved
		if !bytes.Equal(content, testContent) {
			t.Errorf("OpenContent returned different content than what was saved")
		}
	})

	t.Run("Remove", func(t *testing.T) {
		// First verify the file exists
		source, _ := store.Source()
		savedPath := filepath.Join(source, testID)
		if _, err := os.Stat(savedPath); os.IsNotExist(err) {
			t.Fatalf("Test setup failed: file doesn't exist before removal")
		}

		// Remove the content
		err := store.Remove(ctx, testID)
		if err != nil {
			t.Fatalf("Remove() failed: %v", err)
		}

		// Verify the file no longer exists
		if _, err := os.Stat(savedPath); !os.IsNotExist(err) {
			t.Errorf("File still exists after Remove")
		}

		// Try to open the removed content and verify it fails
		_, err = store.Open(ctx, testID)
		if err == nil {
			t.Errorf("OpenContent succeeded for removed content")
		}
	})

	t.Run("ErrorCases", func(t *testing.T) {
		// Test opening non-existent content
		_, err := store.Open(ctx, "nonexistent")
		if err == nil {
			t.Errorf("OpenContent did not fail for non-existent content")
		}

		// Test removing non-existent content
		err = store.Remove(ctx, "nonexistent")
		if err == nil {
			t.Errorf("RemoveContent did not fail for non-existent content")
		}
	})
}
