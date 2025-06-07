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
	t.Run("Load_FileNotExists", func(t *testing.T) {
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
	t.Run("Save_And_Load", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create test overlays
		testOverlays := []*overlay.Overlay{
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
		err = store.Save(ctx, mockService, false)
		if err != nil {
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
			SetOverlays(gomock.Any()).
			DoAndReturn(func(seq iter.Seq2[*overlay.Overlay, error]) error {
				// Collect overlays from the iterator
				var loadedOverlays []*overlay.Overlay
				for ov, err := range seq {
					if err != nil {
						t.Errorf("Unexpected error in overlay iterator: %v", err)
						return nil
					}
					loadedOverlays = append(loadedOverlays, ov)
				}

				// Verify loaded overlays match what we saved
				if len(loadedOverlays) != len(testOverlays) {
					t.Errorf("Load() returned %d overlays, want %d", len(loadedOverlays), len(testOverlays))
					return nil
				}

				// Create copies of the slices and sort them for comparison
				sortedLoaded := make([]*overlay.Overlay, len(loadedOverlays))
				copy(sortedLoaded, loadedOverlays)
				expectedOverlays := make([]*overlay.Overlay, len(testOverlays))
				copy(expectedOverlays, testOverlays)

				sortOverlays(sortedLoaded)
				sortOverlays(expectedOverlays)

				for i := range sortedLoaded {
					if sortedLoaded[i].RepoPattern != expectedOverlays[i].RepoPattern ||
						sortedLoaded[i].ForInit != expectedOverlays[i].ForInit ||
						sortedLoaded[i].RelativePath != expectedOverlays[i].RelativePath ||
						sortedLoaded[i].ContentLocation != expectedOverlays[i].ContentLocation {
						t.Errorf("Overlay %d mismatch:\nGot: %+v\nWant: %+v",
							i, sortedLoaded[i], expectedOverlays[i])
					}
				}
				return nil
			})

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
	t.Run("Save_NoChanges", func(t *testing.T) {
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
	t.Run("Save_Force", func(t *testing.T) {
		// Create a mock controller
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create test overlays
		testOverlays := []*overlay.Overlay{
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
	t.Run("Load_InvalidFile", func(t *testing.T) {
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
func makeOverlayIterator(overlays []*overlay.Overlay) iter.Seq2[*overlay.Overlay, error] {
	return func(yield func(*overlay.Overlay, error) bool) {
		for _, ov := range overlays {
			if !yield(ov, nil) {
				break
			}
		}
	}
}

// Helper function to sort overlays by RepoPattern for consistent comparison
func sortOverlays(overlays []*overlay.Overlay) {
	if overlays == nil {
		return
	}
	for i := range overlays {
		for j := i + 1; j < len(overlays); j++ {
			if overlays[i].RepoPattern > overlays[j].RepoPattern {
				overlays[i], overlays[j] = overlays[j], overlays[i]
			}
		}
	}
}
