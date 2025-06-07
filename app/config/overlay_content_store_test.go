package config_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/overlay"
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

	// Test 1: Test Source() method
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
	testOverlay := overlay.Overlay{
		RepoPattern:     "github.com/kyoh86/gogh",
		ForInit:         false,
		RelativePath:    "test-path",
		ContentLocation: "",
	}
	testContent := []byte("test content data")

	// Test 2: Test SaveContent method
	var location string
	t.Run("SaveContent", func(t *testing.T) {
		buffer := bytes.NewBuffer(testContent)
		var err error
		location, err = store.SaveContent(ctx, testOverlay, buffer)
		if err != nil {
			t.Fatalf("SaveContent() failed: %v", err)
		}

		// Verify location is not empty and looks like a hash
		if len(location) != 64 || !isHexString(location) {
			t.Errorf("SaveContent returned unexpected location: %q", location)
		}

		// Verify the content was saved to the expected path
		source, _ := store.Source()
		savedPath := filepath.Join(source, location)
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

	// Test 3: Test OpenContent method
	t.Run("OpenContent", func(t *testing.T) {
		reader, err := store.OpenContent(ctx, location)
		if err != nil {
			t.Fatalf("OpenContent() failed: %v", err)
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

	// Test 4: Test consistency of content locations
	t.Run("ContentLocationConsistency", func(t *testing.T) {
		// Test that saving the same overlay content twice produces the same location
		buffer1 := bytes.NewBuffer(testContent)
		location1, err := store.SaveContent(ctx, testOverlay, buffer1)
		if err != nil {
			t.Fatalf("First SaveContent() failed: %v", err)
		}

		buffer2 := bytes.NewBuffer(testContent)
		location2, err := store.SaveContent(ctx, testOverlay, buffer2)
		if err != nil {
			t.Fatalf("Second SaveContent() failed: %v", err)
		}

		if location1 != location2 {
			t.Errorf("SaveContent generated different locations for the same overlay and content: %q vs %q",
				location1, location2)
		}

		// Test that different overlays generate different locations
		differentOverlay := overlay.Overlay{
			RepoPattern:     "gitlab.com/kyoh86/gogh",
			ForInit:         false,
			RelativePath:    "test-path",
			ContentLocation: "",
		}

		buffer3 := bytes.NewBuffer(testContent)
		location3, err := store.SaveContent(ctx, differentOverlay, buffer3)
		if err != nil {
			t.Fatalf("SaveContent() failed with different overlay: %v", err)
		}

		if location1 == location3 {
			t.Errorf("SaveContent generated the same location for different overlays")
		}
	})

	// Test 5: Test RemoveContent method
	t.Run("RemoveContent", func(t *testing.T) {
		// First verify the file exists
		source, _ := store.Source()
		savedPath := filepath.Join(source, location)
		if _, err := os.Stat(savedPath); os.IsNotExist(err) {
			t.Fatalf("Test setup failed: file doesn't exist before removal")
		}

		// Remove the content
		err := store.RemoveContent(ctx, location)
		if err != nil {
			t.Fatalf("RemoveContent() failed: %v", err)
		}

		// Verify the file no longer exists
		if _, err := os.Stat(savedPath); !os.IsNotExist(err) {
			t.Errorf("File still exists after RemoveContent")
		}

		// Try to open the removed content and verify it fails
		_, err = store.OpenContent(ctx, location)
		if err == nil {
			t.Errorf("OpenContent succeeded for removed content")
		}
	})

	// Test 6: Error cases
	t.Run("ErrorCases", func(t *testing.T) {
		// Test opening non-existent content
		_, err := store.OpenContent(ctx, "nonexistent")
		if err == nil {
			t.Errorf("OpenContent did not fail for non-existent content")
		}

		// Test removing non-existent content
		err = store.RemoveContent(ctx, "nonexistent")
		if err == nil {
			t.Errorf("RemoveContent did not fail for non-existent content")
		}
	})
}

// Helper function to check if a string contains only hexadecimal characters
func isHexString(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
