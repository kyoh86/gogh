package config_test

import (
	"context"
	"iter"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestOverlayStore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "overlay-store-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock AppContextPathFunc to use our temp directory
	origAppContextPathFunc := config.AppContextPathFunc
	defer func() { config.AppContextPathFunc = origAppContextPathFunc }()

	overlayFile := filepath.Join(tempDir, "overlay.v4.toml")
	config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		if envar == "GOGH_OVERLAY_PATH" {
			return overlayFile, nil
		}
		return "", nil
	}

	// Create test context and store
	ctx := context.Background()
	store := config.NewOverlayStore()

	// Test 1: Test Source() method
	t.Run("Source", func(t *testing.T) {
		source, err := store.Source()
		if err != nil {
			t.Fatalf("Source() failed: %v", err)
		}
		if source != overlayFile {
			t.Errorf("Source() = %q, want %q", source, overlayFile)
		}
	})

	// Test 2: Test Load when file does not exist
	t.Run("Load FileNotExists", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a mock OverlayService as the initial service
		mockInitial := overlay_mock.NewMockOverlayService(ctrl)
		mockInitial.EXPECT().MarkSaved()

		// Test loading when the file doesn't exist
		service, err := store.Load(ctx, func() overlay.OverlayService {
			return mockInitial
		})
		if err != nil {
			t.Fatalf("Load() failed when file doesn't exist: %v", err)
		}
		if service != mockInitial {
			t.Errorf("Load() returned different service than initial")
		}
	})

	// Test 3: Test Save and Load
	t.Run("Save And Load", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create test overlays
		testOverlays := []overlay.Overlay{
			{
				RepoPattern:     "github.com/kyoh86/gogh",
				ForInit:         false,
				RelativePath:    "path1",
				ContentLocation: "location1",
			},
			{
				RepoPattern:     "github.com/kyoh86/another",
				ForInit:         true,
				RelativePath:    "path2",
				ContentLocation: "location2",
			},
		}

		// Create a mock OverlayService for saving
		mockService := overlay_mock.NewMockOverlayService(ctrl)
		mockService.EXPECT().HasChanges().Return(true)
		mockService.EXPECT().ListOverlays().Return(makeOverlayIterator(testOverlays))
		mockService.EXPECT().MarkSaved()

		// Test saving overlays
		if err := store.Save(ctx, mockService, false); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}

		// Verify the file exists
		if _, err := os.Stat(overlayFile); os.IsNotExist(err) {
			t.Fatalf("Save() did not create overlay file")
		}

		// Now test loading the saved overlays
		mockLoadService := overlay_mock.NewMockOverlayService(ctrl)

		// SetOverlays expects an iterator, not a slice
		mockLoadService.EXPECT().
			SetOverlays(testOverlays).
			Return(nil)

		mockLoadService.EXPECT().MarkSaved()

		// Load the overlays
		service, err := store.Load(ctx, func() overlay.OverlayService {
			return mockLoadService
		})
		if err != nil {
			t.Fatalf("Load() failed: %v", err)
		}
		if service != mockLoadService {
			t.Errorf("Load() returned different service than provided")
		}
	})

	// Test 4: Test Save with no changes and force=false
	t.Run("Save NoChanges", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a mock OverlayService
		mockService := overlay_mock.NewMockOverlayService(ctrl)
		mockService.EXPECT().HasChanges().Return(false)

		// Test saving with no changes and force=false
		err = store.Save(ctx, mockService, false)
		if err != nil {
			t.Fatalf("Save() failed with no changes: %v", err)
		}
		// No additional expectations - the function should return early
	})

	// Test 5: Test Save with no changes but force=true
	t.Run("Save Force", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create test overlays
		testOverlays := []overlay.Overlay{
			{
				RepoPattern:     "github.com/kyoh86/forced",
				ForInit:         true,
				RelativePath:    "path3",
				ContentLocation: "location3",
			},
		}

		// Create a mock OverlayService
		mockService := overlay_mock.NewMockOverlayService(ctrl)
		mockService.EXPECT().HasChanges().Return(false) // Force=true still checks HasChanges
		mockService.EXPECT().ListOverlays().Return(makeOverlayIterator(testOverlays))
		mockService.EXPECT().MarkSaved()

		// Test saving with force=true
		err = store.Save(ctx, mockService, true)
		if err != nil {
			t.Fatalf("Save() failed with force=true: %v", err)
		}

		// Verify the file exists and was updated
		fileInfo, err := os.Stat(overlayFile)
		if os.IsNotExist(err) {
			t.Fatalf("Save() did not create overlay file with force=true")
		}
		if fileInfo.Size() == 0 {
			t.Errorf("Save() created empty file with force=true")
		}
	})

	// Test 6: Test Load with invalid file content
	t.Run("Load InvalidFile", func(t *testing.T) {
		// Create an invalid TOML file
		invalidContent := `
		This is not valid TOML content
		overlays = [
		`
		if err := os.WriteFile(overlayFile, []byte(invalidContent), 0644); err != nil {
			t.Fatalf("Failed to write invalid overlay file: %v", err)
		}

		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create a mock OverlayService
		mockService := overlay_mock.NewMockOverlayService(ctrl)

		// Attempt to load the invalid file
		_, err := store.Load(ctx, func() overlay.OverlayService {
			return mockService
		})
		if err == nil {
			t.Errorf("Load() did not fail with invalid file content")
		}
	})
}

// Helper function to create an iterator for overlays
func makeOverlayIterator(overlays []overlay.Overlay) iter.Seq2[*overlay.Overlay, error] {
	return func(yield func(*overlay.Overlay, error) bool) {
		for _, ov := range overlays {
			ov := ov
			if !yield(&ov, nil) {
				break
			}
		}
	}
}
