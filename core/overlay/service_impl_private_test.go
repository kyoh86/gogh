package overlay

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

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

	if impl.contentStore != store {
		t.Error("expected contentStore to be set correctly")
	}

	if len(impl.overlays) != 0 {
		t.Errorf("expected empty overlays, got %d", len(impl.overlays))
	}

	if impl.changed {
		t.Error("expected changed to be false initially")
	}
}

func TestListOverlays(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)

	// Empty list
	overlays := service.ListOverlays()
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
	testOverlays := []Overlay{
		{RepoPattern: "repo1", RelativePath: "path1", ForInit: true},
		{RepoPattern: "repo2", RelativePath: "path2", ForInit: false},
	}

	for _, ov := range testOverlays {
		err := service.AddOverlay(context.Background(), ov, strings.NewReader("content"))
		if err != nil {
			t.Fatalf("failed to add overlay: %v", err)
		}
	}

	// Check list again
	resultOverlays, err := typ.CollectWithError(service.ListOverlays())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(resultOverlays) != len(testOverlays) {
		t.Errorf("expected %d overlays, got %d", len(testOverlays), len(resultOverlays))
	}

	// Verify overlays content
	for i, ov := range resultOverlays {
		if ov.RepoPattern != testOverlays[i].RepoPattern ||
			ov.RelativePath != testOverlays[i].RelativePath ||
			ov.ForInit != testOverlays[i].ForInit {
			t.Errorf("overlay mismatch at index %d: expected %v, got %v", i, testOverlays[i], ov)
		}
	}
}

func TestAddOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Test successful add
	ov := Overlay{RepoPattern: "repo1", RelativePath: "path1", ForInit: true}
	content := "test content"
	err := service.AddOverlay(ctx, ov, strings.NewReader(content))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(service.overlays) != 1 {
		t.Errorf("expected 1 overlay, got %d", len(service.overlays))
	}

	if !service.changed {
		t.Error("expected changed to be true after add")
	}

	// Verify the content location was set
	if service.overlays[0].ContentLocation == "" {
		t.Error("expected ContentLocation to be set")
	}

	// Verify the content was actually saved to the filesystem
	savedLocation := service.overlays[0].ContentLocation
	reader, err := store.OpenContent(ctx, savedLocation)
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
	err = service.AddOverlay(ctx, ov, nil)
	if err == nil {
		t.Error("expected error when adding nil content")
	}
}

func TestRemoveOverlay(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add test overlays
	overlays := []Overlay{
		{RepoPattern: "repo1", RelativePath: "path1", ForInit: true},
		{RepoPattern: "repo2", RelativePath: "path2", ForInit: false},
		{RepoPattern: "repo3", RelativePath: "path3", ForInit: true},
	}

	locations := make([]string, 0, len(overlays))
	for _, ov := range overlays {
		err := service.AddOverlay(ctx, ov, strings.NewReader("content"))
		if err != nil {
			t.Fatalf("failed to add overlay: %v", err)
		}
		locations = append(locations, service.overlays[len(service.overlays)-1].ContentLocation)
	}

	service.changed = false // Reset for testing

	// Remove existing overlay (the middle one)
	err := service.RemoveOverlay(ctx, overlays[1])
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(service.overlays) != 2 {
		t.Errorf("expected 2 overlays, got %d", len(service.overlays))
	}

	if !service.changed {
		t.Error("expected changed to be true after remove")
	}

	// Check that the remaining overlays are the expected ones
	if service.overlays[0].RepoPattern != "repo1" || service.overlays[1].RepoPattern != "repo3" {
		t.Error("wrong overlay was removed")
	}

	// Verify content was removed from filesystem
	removedLocation := locations[1]
	_, err = store.OpenContent(ctx, removedLocation)
	if err == nil {
		t.Error("expected error opening removed content, but got none")
	}

	// Remove non-existent overlay
	nonExistent := Overlay{RepoPattern: "non", RelativePath: "existent", ForInit: false}
	service.changed = false // Reset for testing
	err = service.RemoveOverlay(ctx, nonExistent)
	if err != nil {
		t.Errorf("unexpected error removing non-existent overlay: %v", err)
	}

	if len(service.overlays) != 2 {
		t.Errorf("expected still 2 overlays, got %d", len(service.overlays))
	}

	if service.changed {
		t.Error("expected changed to be false when removing non-existent overlay")
	}
}

func TestOpenOverlayContent(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Add test overlay
	ov := Overlay{RepoPattern: "repo1", RelativePath: "path1", ForInit: true}
	expectedContent := "test content"
	err := service.AddOverlay(ctx, ov, strings.NewReader(expectedContent))
	if err != nil {
		t.Fatalf("failed to add overlay: %v", err)
	}

	// Open existing overlay
	reader, err := service.OpenOverlayContent(ctx, ov)
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
	nonExistent := Overlay{RepoPattern: "non", RelativePath: "existent", ForInit: false}
	reader, err = service.OpenOverlayContent(ctx, nonExistent)
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
	ov := Overlay{RepoPattern: "repo", RelativePath: "path", ForInit: true}
	err := service.AddOverlay(ctx, ov, strings.NewReader("content"))
	if err != nil {
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

func TestSetOverlays(t *testing.T) {
	store := NewMockContentStore()
	service := NewOverlayService(store).(*serviceImpl)
	ctx := context.Background()

	// Create test overlays and save content for them
	testOverlays := []*Overlay{
		{RepoPattern: "repo1", RelativePath: "path1", ForInit: true},
		{RepoPattern: "repo2", RelativePath: "path2", ForInit: false},
	}

	// Save content for each overlay
	for _, ov := range testOverlays {
		location, err := store.SaveContent(ctx, *ov, strings.NewReader("content for "+ov.RepoPattern))
		if err != nil {
			t.Fatalf("failed to save content: %v", err)
		}
		ov.ContentLocation = location
	}

	// Create iterator function
	iter := func(yield func(*Overlay, error) bool) {
		for _, ov := range testOverlays {
			if !yield(ov, nil) {
				return
			}
		}
	}

	// Set overlays
	err := service.SetOverlays(iter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(service.overlays) != len(testOverlays) {
		t.Errorf("expected %d overlays, got %d", len(testOverlays), len(service.overlays))
	}

	for i, ov := range service.overlays {
		if !reflect.DeepEqual(ov, *testOverlays[i]) {
			t.Errorf("overlay mismatch at index %d: expected %v, got %v", i, *testOverlays[i], ov)
		}
	}

	if service.changed {
		t.Error("expected changed to be false after SetOverlays")
	}

	// Test with error in iterator
	expectedError := errors.New("test error")
	errIter := func(yield func(*Overlay, error) bool) {
		yield(nil, expectedError)
	}

	err = service.SetOverlays(errIter)
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
