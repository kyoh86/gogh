package overlay

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Test Update method thoroughly
func TestUpdateOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add initial overlay
	initialContent := "initial content"
	initial := Entry{
		Name:         "test-overlay",
		RelativePath: "path/to/file.txt",
		Content:      strings.NewReader(initialContent),
	}
	id, err := service.Add(ctx, initial)
	if err != nil {
		t.Fatalf("failed to add overlay: %v", err)
	}

	service.MarkSaved() // Reset dirty flag

	testCases := []struct {
		name      string
		update    Entry
		wantDirty bool
		check     func(t *testing.T, ov Overlay)
	}{
		{
			name: "update content only",
			update: Entry{
				Content: strings.NewReader("updated content"),
			},
			wantDirty: true,
			check: func(t *testing.T, ov Overlay) {
				// Content should be updated
				reader, err := store.Open(ctx, ov.ID())
				if err != nil {
					t.Fatalf("failed to open updated content: %v", err)
				}
				defer reader.Close()
				content, _ := io.ReadAll(reader)
				if string(content) != "updated content" {
					t.Errorf("content not updated: got %q", string(content))
				}
				// Name and path should remain unchanged
				if ov.Name() != "test-overlay" {
					t.Errorf("name changed unexpectedly: got %q", ov.Name())
				}
				if ov.RelativePath() != "path/to/file.txt" {
					t.Errorf("path changed unexpectedly: got %q", ov.RelativePath())
				}
			},
		},
		{
			name: "update name only",
			update: Entry{
				Name: "renamed-overlay",
			},
			wantDirty: true,
			check: func(t *testing.T, ov Overlay) {
				if ov.Name() != "renamed-overlay" {
					t.Errorf("name not updated: got %q", ov.Name())
				}
			},
		},
		{
			name: "update path only",
			update: Entry{
				RelativePath: "new/path/file.txt",
			},
			wantDirty: true,
			check: func(t *testing.T, ov Overlay) {
				if ov.RelativePath() != "new/path/file.txt" {
					t.Errorf("path not updated: got %q", ov.RelativePath())
				}
			},
		},
		{
			name: "update all fields",
			update: Entry{
				Name:         "final-name",
				RelativePath: "final/path.txt",
				Content:      strings.NewReader("final content"),
			},
			wantDirty: true,
			check: func(t *testing.T, ov Overlay) {
				if ov.Name() != "final-name" {
					t.Errorf("name not updated: got %q", ov.Name())
				}
				if ov.RelativePath() != "final/path.txt" {
					t.Errorf("path not updated: got %q", ov.RelativePath())
				}
			},
		},
		{
			name:      "update with empty entry (no changes)",
			update:    Entry{},
			wantDirty: false,
			check: func(t *testing.T, ov Overlay) {
				// Nothing should change
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service.MarkSaved() // Reset before each test

			err := service.Update(ctx, id, tc.update)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if service.dirty != tc.wantDirty {
				t.Errorf("dirty flag: got %v, want %v", service.dirty, tc.wantDirty)
			}

			// Get the updated overlay
			ov, err := service.Get(ctx, id)
			if err != nil {
				t.Fatalf("failed to get overlay: %v", err)
			}

			tc.check(t, ov)
		})
	}

	// Test update non-existent overlay
	err = service.Update(ctx, "non-existent", Entry{Name: "new"})
	if err == nil {
		t.Error("expected error when updating non-existent overlay")
	}

	// Test update with partial ID
	if len(id) > 8 {
		partialID := id[:8]
		err = service.Update(ctx, partialID, Entry{Name: "partial-update"})
		if err != nil {
			t.Errorf("failed to update with partial ID: %v", err)
		}

		ov, _ := service.Get(ctx, id)
		if ov.Name() != "partial-update" {
			t.Error("partial ID update failed")
		}
	}

	// Test content save error (would require a mock that can simulate errors)
	// This is currently not testable with the simple MockContentStore
	// In a real implementation, you would use a mock that can be configured to return errors
}

// Test Get method thoroughly
func TestGetOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add test overlays
	overlay1 := Entry{
		Name:         "overlay1",
		RelativePath: "path1.txt",
		Content:      strings.NewReader("content1"),
	}
	id1, _ := service.Add(ctx, overlay1)

	overlay2 := Entry{
		Name:         "overlay2",
		RelativePath: "path2.txt",
		Content:      strings.NewReader("content2"),
	}
	id2, _ := service.Add(ctx, overlay2)

	testCases := []struct {
		name    string
		idlike  string
		wantID  string
		wantErr bool
	}{
		{
			name:    "get by full ID",
			idlike:  id1,
			wantID:  id1,
			wantErr: false,
		},
		{
			name:    "get by partial ID",
			idlike:  id1[:8],
			wantID:  id1,
			wantErr: false,
		},
		{
			name:    "get non-existent",
			idlike:  "non-existent",
			wantErr: true,
		},
		{
			name:    "get with empty ID",
			idlike:  "",
			wantErr: true,
		},
		{
			name:    "get second overlay",
			idlike:  id2,
			wantID:  id2,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ov, err := service.Get(ctx, tc.idlike)
			if (err != nil) != tc.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && ov.ID() != tc.wantID {
				t.Errorf("Get() returned wrong overlay: got %s, want %s", ov.ID(), tc.wantID)
			}
		})
	}
}

// Test Load method
func TestLoadOverlays(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add some overlays first
	initial := Entry{
		Name:         "initial",
		RelativePath: "initial.txt",
		Content:      strings.NewReader("initial"),
	}
	_, _ = service.Add(ctx, initial)
	service.MarkSaved()

	// Create overlays to load
	overlaysToLoad := []Overlay{
		&overlayElement{
			id:           uuid.New(),
			name:         "loaded1",
			relativePath: "loaded1.txt",
		},
		&overlayElement{
			id:           uuid.New(),
			name:         "loaded2",
			relativePath: "loaded2.txt",
		},
		&overlayElement{
			id:           uuid.New(),
			name:         "loaded3",
			relativePath: "loaded3.txt",
		},
	}

	// Load new overlays
	loadSeq := func(yield func(Overlay, error) bool) {
		for _, ov := range overlaysToLoad {
			if !yield(ov, nil) {
				return
			}
		}
	}

	err := service.Load(loadSeq)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check that old overlays are replaced
	if service.overlays.Len() != len(overlaysToLoad) {
		t.Errorf("expected %d overlays after load, got %d", len(overlaysToLoad), service.overlays.Len())
	}

	// Verify loaded overlays
	for _, expected := range overlaysToLoad {
		ov, err := service.Get(ctx, expected.ID())
		if err != nil {
			t.Errorf("failed to get loaded overlay %s: %v", expected.ID(), err)
		}
		if ov.Name() != expected.Name() {
			t.Errorf("name mismatch: got %q, want %q", ov.Name(), expected.Name())
		}
		if ov.RelativePath() != expected.RelativePath() {
			t.Errorf("path mismatch: got %q, want %q", ov.RelativePath(), expected.RelativePath())
		}
	}

	// Check dirty flag
	if !service.dirty {
		t.Error("expected dirty to be true after load")
	}

	// Test load with error
	errorSeq := func(yield func(Overlay, error) bool) {
		yield(nil, errors.New("load error"))
	}

	err = service.Load(errorSeq)
	if err == nil {
		t.Error("expected error from load")
	}

	// Test partial load with error
	partialSeq := func(yield func(Overlay, error) bool) {
		yield(overlaysToLoad[0], nil)
		yield(nil, errors.New("partial load error"))
	}

	service.overlays.Clear() // Clear before partial load test
	err = service.Load(partialSeq)
	if err == nil {
		t.Error("expected error from partial load")
	}
}

// Test concurrent operations
func TestConcurrentOperations(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add initial overlays
	var ids []string
	for i := 0; i < 10; i++ {
		id, err := service.Add(ctx, Entry{
			Name:         "overlay" + string(rune('0'+i)),
			RelativePath: "path" + string(rune('0'+i)) + ".txt",
			Content:      strings.NewReader("content" + string(rune('0'+i))),
		})
		if err != nil {
			t.Fatalf("failed to add overlay: %v", err)
		}
		ids = append(ids, id)
	}

	// Run concurrent operations
	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Reader goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				// List
				count := 0
				for _, err := range service.List() {
					if err != nil {
						errChan <- err
						return
					}
					count++
				}

				// Get random overlay
				idx := j % len(ids)
				_, err := service.Get(ctx, ids[idx])
				if err != nil {
					errChan <- err
					return
				}

				// Open random overlay
				_, err = service.Open(ctx, ids[idx])
				if err != nil {
					errChan <- err
					return
				}

				time.Sleep(time.Microsecond)
			}
		}()
	}

	// Writer goroutines
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// Update
				idx := j % len(ids)
				err := service.Update(ctx, ids[idx], Entry{
					Name: "updated-" + string(rune('0'+workerID)) + "-" + string(rune('0'+j)),
				})
				if err != nil {
					errChan <- err
					return
				}

				// Add new
				_, err = service.Add(ctx, Entry{
					Name:         "new-" + string(rune('0'+workerID)) + "-" + string(rune('0'+j)),
					RelativePath: "new-path.txt",
					Content:      strings.NewReader("new content"),
				})
				if err != nil {
					errChan <- err
					return
				}

				time.Sleep(time.Microsecond * 10)
			}
		}(i)
	}

	// Wait for all goroutines
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("concurrent operation error: %v", err)
	}
}

// Test edge cases
func TestOverlayServiceEdgeCases(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	t.Run("add with large content", func(t *testing.T) {
		// Create large content
		var buf bytes.Buffer
		for i := 0; i < 1000000; i++ {
			buf.WriteByte(byte(i % 256))
		}

		id, err := service.Add(ctx, Entry{
			Name:         "large",
			RelativePath: "large.bin",
			Content:      &buf,
		})
		if err != nil {
			t.Fatalf("failed to add large overlay: %v", err)
		}

		// Verify content
		reader, err := service.Open(ctx, id)
		if err != nil {
			t.Fatalf("failed to open large overlay: %v", err)
		}
		defer reader.Close()

		content, _ := io.ReadAll(reader)
		if len(content) != 1000000 {
			t.Errorf("content size mismatch: got %d, want 1000000", len(content))
		}
	})

	t.Run("add with empty name and path", func(t *testing.T) {
		id, err := service.Add(ctx, Entry{
			Name:         "",
			RelativePath: "",
			Content:      strings.NewReader("content"),
		})
		if err != nil {
			t.Fatalf("failed to add overlay with empty fields: %v", err)
		}

		ov, _ := service.Get(ctx, id)
		if ov.Name() != "" {
			t.Errorf("expected empty name, got %q", ov.Name())
		}
		if ov.RelativePath() != "" {
			t.Errorf("expected empty path, got %q", ov.RelativePath())
		}
	})

	t.Run("update to empty values", func(t *testing.T) {
		id, _ := service.Add(ctx, Entry{
			Name:         "test",
			RelativePath: "test.txt",
			Content:      strings.NewReader("test"),
		})

		// Update with empty strings (should keep existing values)
		err := service.Update(ctx, id, Entry{
			Name:         "",
			RelativePath: "",
		})
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}

		ov, _ := service.Get(ctx, id)
		if ov.Name() != "test" {
			t.Error("name was cleared unexpectedly")
		}
		if ov.RelativePath() != "test.txt" {
			t.Error("path was cleared unexpectedly")
		}
	})

	t.Run("remove during iteration", func(t *testing.T) {
		// Add overlays
		var ids []string
		for i := 0; i < 5; i++ {
			id, _ := service.Add(ctx, Entry{
				Name:         "iter" + string(rune('0'+i)),
				RelativePath: "path.txt",
				Content:      strings.NewReader("content"),
			})
			ids = append(ids, id)
		}

		// Try to remove during iteration (should be safe due to mutex)
		count := 0
		for ov, err := range service.List() {
			if err != nil {
				t.Errorf("list error: %v", err)
			}
			if count == 2 {
				// Remove an overlay during iteration
				go func() {
					_ = service.Remove(ctx, ids[3])
				}()
			}
			_ = ov
			count++
		}
	})
}

// Test HasChanges and MarkSaved with all operations
func TestDirtyFlagManagement(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Initially clean
	if service.HasChanges() {
		t.Error("expected no changes initially")
	}

	// Add makes it dirty
	id, _ := service.Add(ctx, Entry{
		Name:         "test",
		RelativePath: "test.txt",
		Content:      strings.NewReader("test"),
	})
	if !service.HasChanges() {
		t.Error("expected changes after add")
	}

	// MarkSaved cleans it
	service.MarkSaved()
	if service.HasChanges() {
		t.Error("expected no changes after MarkSaved")
	}

	// Update makes it dirty
	_ = service.Update(ctx, id, Entry{Name: "updated"})
	if !service.HasChanges() {
		t.Error("expected changes after update")
	}

	service.MarkSaved()

	// Remove makes it dirty
	_ = service.Remove(ctx, id)
	if !service.HasChanges() {
		t.Error("expected changes after remove")
	}

	service.MarkSaved()

	// Load makes it dirty
	loadSeq := func(yield func(Overlay, error) bool) {
		ov := &overlayElement{
			id:           uuid.New(),
			name:         "loaded",
			relativePath: "loaded.txt",
		}
		yield(ov, nil)
	}
	_ = service.Load(loadSeq)
	if !service.HasChanges() {
		t.Error("expected changes after load")
	}

	// Failed operations should not make it dirty
	service.MarkSaved()

	// Failed update
	_ = service.Update(ctx, "non-existent", Entry{Name: "fail"})
	if service.HasChanges() {
		t.Error("expected no changes after failed update")
	}

	// Failed remove
	_ = service.Remove(ctx, "non-existent")
	if service.HasChanges() {
		t.Error("expected no changes after failed remove")
	}
}
