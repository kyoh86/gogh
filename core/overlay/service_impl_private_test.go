package overlay

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/typ"
)

func TestNewOverlayService(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store)

	if service == nil {
		t.Fatal("expected service to not be nil")
	}

	impl, ok := service.(*serviceImpl)
	if !ok {
		t.Fatal("expected service to be a *serviceImpl")
	}

	if impl.content != store {
		t.Error("expected contentStore to be set correctly")
	}

	if impl.overlays.Len() != 0 {
		t.Errorf("expected empty overlays, got %d", impl.overlays.Len())
	}

	if impl.dirty {
		t.Error("expected dirty to be false initially")
	}
}

func TestListOverlays(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)

	// Empty list
	overlays := service.List()
	count := 0
	for ov, err := range overlays {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		count++
		_ = ov // Avoid unused variable warning
	}
	if count != 0 {
		t.Errorf("expected 0 overlays, got %d", count)
	}

	// Add some overlays
	testOverlays := []Entry{
		{RelativePath: "path1", Content: strings.NewReader("content1")},
		{RelativePath: "path2", Content: strings.NewReader("content2")},
	}

	for _, ov := range testOverlays {
		if _, err := service.Add(context.Background(), ov); err != nil {
			t.Fatalf("failed to add overlay: %v", err)
		}
	}

	// Check list again
	resultOverlays, err := typ.CollectWithError(service.List())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(resultOverlays) != len(testOverlays) {
		t.Errorf("expected %d overlays, got %d", len(testOverlays), len(resultOverlays))
	}

	// Verify overlays content
	for i, ov := range resultOverlays {
		if ov.RelativePath() != testOverlays[i].RelativePath {
			t.Errorf("overlay mismatch at index %d: expected %v, got %v", i, testOverlays[i], ov)
		}
	}
}

func TestAddOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Test successful add
	ov := Entry{RelativePath: "path1", Content: strings.NewReader("test content")}
	content := "test content"
	id, err := service.Add(ctx, ov)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.overlays.Len() != 1 {
		t.Errorf("expected 1 overlay, got %d", service.overlays.Len())
	}

	if !service.dirty {
		t.Error("expected dirty to be true after add")
	}

	// Verify the content was actually saved to the filesystem
	reader, err := store.Open(ctx, id)
	if err != nil {
		t.Errorf("failed to open saved content: %v", err)
	}
	defer reader.Close()

	savedContent, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("failed to read saved content: %v", err)
	}

	if string(savedContent) != content {
		t.Errorf("content mismatch: expected %q, got %q", content, string(savedContent))
	}

	// Test add with nil content
	if _, err := service.Add(ctx, Entry{RelativePath: ov.RelativePath}); err == nil {
		t.Error("expected error when adding nil content")
	}

	// Test adding duplicate overlay
	dupID, err := service.Add(ctx, ov)
	if err != nil {
		t.Errorf("unexpected error adding duplicate overlay: %v", err)
	}
	if dupID == id {
		t.Error("expected different ID for duplicate overlay, got same ID")
	}
}

func TestRemoveOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add test overlays
	overlays := []Entry{
		{RelativePath: "path1", Content: strings.NewReader("content1")},
		{RelativePath: "path2", Content: strings.NewReader("content2")},
		{RelativePath: "path3", Content: strings.NewReader("content3")},
	}

	ids := make([]string, 0, len(overlays))
	for _, ov := range overlays {
		id, err := service.Add(ctx, ov)
		ids = append(ids, id)
		if err != nil {
			t.Fatalf("failed to add overlay: %v", err)
		}
	}

	service.dirty = false // Reset for testing

	// Remove existing overlay (the middle one)
	err := service.Remove(ctx, ids[1])
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.overlays.Len() != 2 {
		t.Errorf("expected 2 overlays, got %d", service.overlays.Len())
	}

	if !service.dirty {
		t.Error("expected dirty to be true after remove")
	}

	// Check that the remaining overlays are the expected ones
	if service.overlays.At(0).relativePath != "path1" || service.overlays.At(1).relativePath != "path3" {
		t.Error(service.overlays.At(0), service.overlays.At(1))
		t.Error("wrong overlay was removed")
	}

	// Verify content was removed from filesystem
	removedID := ids[1]
	_, err = store.Open(ctx, removedID)
	if err == nil {
		t.Error("expected error opening removed content, but got none")
	}

	// Remove non-existent overlay
	nonExistent := uuid.NewString()
	service.dirty = false // Reset for testing
	err = service.Remove(ctx, nonExistent)
	if err == nil {
		t.Error("expected error when removing non-existent overlay, got nil")
	}

	if service.overlays.Len() != 2 {
		t.Errorf("expected still 2 overlays, got %d", service.overlays.Len())
	}

	if service.dirty {
		t.Error("expected dirty to be false when removing non-existent overlay")
	}
}

func TestOpenOverlayContent(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add test overlay
	expectedContent := "test content"
	ov := Entry{RelativePath: "path1", Content: strings.NewReader(expectedContent)}
	id, err := service.Add(ctx, ov)
	if err != nil {
		t.Fatalf("failed to add overlay: %v", err)
	}

	// Open existing overlay
	reader, err := service.Open(ctx, id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("failed to read content: %v", err)
	}

	if string(content) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(content))
	}

	// Open non-existent overlay
	nonExistent := uuid.NewString()
	reader, err = service.Open(ctx, nonExistent)
	if err == nil {
		t.Error("expected error for non-existent overlay, got nil")
		if reader != nil {
			reader.Close()
		}
	}
}

func TestHasChanges(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Initially no changes
	if service.HasChanges() {
		t.Error("expected no changes initially")
	}

	// After adding overlay
	ov := Entry{RelativePath: "path", Content: strings.NewReader("content")}
	if _, err := service.Add(ctx, ov); err != nil {
		t.Fatalf("failed to add overlay: %v", err)
	}

	if !service.HasChanges() {
		t.Error("expected changes after adding overlay")
	}

	// After marking saved
	service.MarkSaved()
	if service.HasChanges() {
		t.Error("expected no changes after marking saved")
	}
}
