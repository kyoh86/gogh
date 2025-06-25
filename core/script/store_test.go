package script_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/kyoh86/gogh/v4/core/script"
)

// Ensure mockScriptSourceStore implements the interface
var _ script.ScriptSourceStore = (*mockScriptSourceStore)(nil)

func TestScriptSourceStore_Save(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()

	t.Run("save new content", func(t *testing.T) {
		content := bytes.NewReader([]byte("test content"))
		err := store.Save(ctx, "script1", content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify content was saved
		saved, exists := store.contents["script1"]
		if !exists {
			t.Error("expected content to be saved")
		}
		if string(saved) != "test content" {
			t.Errorf("expected 'test content', got %s", string(saved))
		}
	})

	t.Run("overwrite existing content", func(t *testing.T) {
		// Save initial content
		_ = store.Save(ctx, "script2", bytes.NewReader([]byte("initial")))

		// Overwrite with new content
		err := store.Save(ctx, "script2", bytes.NewReader([]byte("updated")))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		saved := store.contents["script2"]
		if string(saved) != "updated" {
			t.Errorf("expected 'updated', got %s", string(saved))
		}
	})

	t.Run("save empty content", func(t *testing.T) {
		err := store.Save(ctx, "script3", bytes.NewReader([]byte("")))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		saved := store.contents["script3"]
		if len(saved) != 0 {
			t.Errorf("expected empty content, got %d bytes", len(saved))
		}
	})

	t.Run("save error", func(t *testing.T) {
		store.saveErr = errors.New("save failed")
		err := store.Save(ctx, "script4", bytes.NewReader([]byte("content")))
		if err == nil {
			t.Error("expected error, got nil")
		}
		store.saveErr = nil
	})
}

func TestScriptSourceStore_Open(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()

	// Pre-populate some content
	store.contents["existing"] = []byte("existing content")
	store.contents["empty"] = []byte("")

	t.Run("open existing content", func(t *testing.T) {
		reader, err := store.Open(ctx, "existing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if string(data) != "existing content" {
			t.Errorf("expected 'existing content', got %s", string(data))
		}
	})

	t.Run("open empty content", func(t *testing.T) {
		reader, err := store.Open(ctx, "empty")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}
		if len(data) != 0 {
			t.Errorf("expected empty content, got %d bytes", len(data))
		}
	})

	t.Run("open non-existent", func(t *testing.T) {
		_, err := store.Open(ctx, "non-existent")
		if err == nil {
			t.Error("expected error when opening non-existent content")
		}
	})

	t.Run("open error", func(t *testing.T) {
		store.openErr = errors.New("open failed")
		_, err := store.Open(ctx, "existing")
		if err == nil {
			t.Error("expected error, got nil")
		}
		store.openErr = nil
	})

	t.Run("reader can be closed", func(t *testing.T) {
		reader, err := store.Open(ctx, "existing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be able to close without error
		err = reader.Close()
		if err != nil {
			t.Errorf("unexpected error closing reader: %v", err)
		}
	})
}

func TestScriptSourceStore_Remove(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()

	// Pre-populate some content
	store.contents["toremove1"] = []byte("content1")
	store.contents["toremove2"] = []byte("content2")
	store.contents["tokeep"] = []byte("keep this")

	t.Run("remove existing content", func(t *testing.T) {
		err := store.Remove(ctx, "toremove1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify content was removed
		_, exists := store.contents["toremove1"]
		if exists {
			t.Error("expected content to be removed")
		}

		// Verify other content still exists
		_, exists = store.contents["tokeep"]
		if !exists {
			t.Error("expected other content to remain")
		}
	})

	t.Run("remove non-existent", func(t *testing.T) {
		// Should not error when removing non-existent
		err := store.Remove(ctx, "non-existent")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("remove error", func(t *testing.T) {
		store.removeErr = errors.New("remove failed")
		err := store.Remove(ctx, "toremove2")
		if err == nil {
			t.Error("expected error, got nil")
		}

		// Content should still exist since remove failed
		_, exists := store.contents["toremove2"]
		if !exists {
			t.Error("content should not be removed when error occurs")
		}
		store.removeErr = nil
	})
}

func TestScriptSourceStore_Integration(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()

	// Save -> Open -> Remove workflow
	t.Run("full workflow", func(t *testing.T) {
		scriptID := "workflow-test"
		content := "workflow content"

		// Save
		err := store.Save(ctx, scriptID, bytes.NewReader([]byte(content)))
		if err != nil {
			t.Fatalf("save failed: %v", err)
		}

		// Open and verify
		reader, err := store.Open(ctx, scriptID)
		if err != nil {
			t.Fatalf("open failed: %v", err)
		}

		data, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if string(data) != content {
			t.Errorf("expected %s, got %s", content, string(data))
		}

		// Remove
		err = store.Remove(ctx, scriptID)
		if err != nil {
			t.Fatalf("remove failed: %v", err)
		}

		// Verify it's gone
		_, err = store.Open(ctx, scriptID)
		if err == nil {
			t.Error("expected error when opening removed content")
		}
	})

	t.Run("multiple scripts", func(t *testing.T) {
		// Save multiple scripts
		scripts := map[string]string{
			"script1": "content1",
			"script2": "content2",
			"script3": "content3",
		}

		for id, content := range scripts {
			err := store.Save(ctx, id, bytes.NewReader([]byte(content)))
			if err != nil {
				t.Fatalf("failed to save %s: %v", id, err)
			}
		}

		// Verify all can be opened
		for id, expectedContent := range scripts {
			reader, err := store.Open(ctx, id)
			if err != nil {
				t.Fatalf("failed to open %s: %v", id, err)
			}

			data, _ := io.ReadAll(reader)
			reader.Close()

			if string(data) != expectedContent {
				t.Errorf("script %s: expected %s, got %s", id, expectedContent, string(data))
			}
		}

		// Remove one
		_ = store.Remove(ctx, "script2")

		// Verify only the removed one is gone
		_, err := store.Open(ctx, "script2")
		if err == nil {
			t.Error("script2 should be removed")
		}

		// Others should still exist
		for _, id := range []string{"script1", "script3"} {
			_, err := store.Open(ctx, id)
			if err != nil {
				t.Errorf("script %s should still exist", id)
			}
		}
	})
}
