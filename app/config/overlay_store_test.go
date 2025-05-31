package config_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/mock/gomock"
)

// setupOverlayStoreTestTempDir creates a temporary directory for testing and returns its path
// along with a cleanup function.
func setupOverlayStoreTestTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "overlay-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// setupOverlayStoreTestEnvironment sets up a test environment with temporary directory
// and mocked OverlayService
func setupOverlayStoreTestEnvironment(t *testing.T) (
	string,
	func(),
	*workspace_mock.MockOverlayService,
	*config.OverlayStore,
) {
	ctrl := gomock.NewController(t)
	tempDir, cleanup := setupOverlayStoreTestTempDir(t)
	mockService := workspace_mock.NewMockOverlayService(ctrl)

	// Override the appContextPath to use our test directory
	origAppContextPathFunc := config.AppContextPathFunc
	config.AppContextPathFunc = func(envName string, fallbackFunc func() (string, error), rel ...string) (string, error) {
		return filepath.Join(append([]string{tempDir}, rel...)...), nil
	}

	// Add cleanup for the override
	originalCleanup := cleanup
	cleanup = func() {
		originalCleanup()
		config.AppContextPathFunc = origAppContextPathFunc
		ctrl.Finish()
	}

	store := config.NewOverlayStore()

	return tempDir, cleanup, mockService, store
}

// createOverlayStoreTestTOMLFile creates a test TOML file with overlay settings
func createOverlayStoreTestTOMLFile(t *testing.T, path string, patterns []workspace.OverlayPattern) {
	t.Helper()

	type tomlOverlayFile struct {
		SourcePath string `toml:"source_path"`
		TargetPath string `toml:"target_path"`
	}

	type tomlOverlayPattern struct {
		Pattern string            `toml:"pattern"`
		Files   []tomlOverlayFile `toml:"files"`
	}

	type tomlOverlayStore struct {
		Patterns []tomlOverlayPattern `toml:"patterns"`
	}

	v := tomlOverlayStore{
		Patterns: make([]tomlOverlayPattern, 0, len(patterns)),
	}

	for _, pattern := range patterns {
		files := make([]tomlOverlayFile, 0, len(pattern.Files))
		for _, file := range pattern.Files {
			files = append(files, tomlOverlayFile{
				SourcePath: file.SourcePath,
				TargetPath: file.TargetPath,
			})
		}
		v.Patterns = append(v.Patterns, tomlOverlayPattern{
			Pattern: pattern.Pattern,
			Files:   files,
		})
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory for TOML file: %v", err)
	}
	encoded, err := toml.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to encode TOML: %v", err)
	}
	if err := os.WriteFile(path, encoded, 0644); err != nil {
		t.Fatalf("Failed to write TOML file: %v", err)
	}
}

func TestNewOverlayStore(t *testing.T) {
	store := config.NewOverlayStore()
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestOverlayStore_Source(t *testing.T) {
	tempDir, cleanup, _, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "overlay.v4.toml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestOverlayStore_Load_FileExists(t *testing.T) {
	tempDir, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Create test patterns
	pattern1 := workspace.OverlayPattern{
		Pattern: "*.go",
		Files: []workspace.OverlayFile{
			{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
			{SourcePath: "/src/file2.go", TargetPath: "/dest/file2.go"},
		},
	}
	pattern2 := workspace.OverlayPattern{
		Pattern: "*.md",
		Files: []workspace.OverlayFile{
			{SourcePath: "/src/file.md", TargetPath: "/dest/file.md"},
		},
	}
	patterns := []workspace.OverlayPattern{pattern1, pattern2}

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "overlay.v4.toml")
	createOverlayStoreTestTOMLFile(t, configPath, patterns)

	// Setup mock expectations
	mockService.EXPECT().AddPattern("*.go", []workspace.OverlayFile{
		{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
		{SourcePath: "/src/file2.go", TargetPath: "/dest/file2.go"},
	}).Return(nil)
	mockService.EXPECT().AddPattern("*.md", []workspace.OverlayFile{
		{SourcePath: "/src/file.md", TargetPath: "/dest/file.md"},
	}).Return(nil)
	mockService.EXPECT().MarkSaved()

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.OverlayService {
		return mockService
	}

	service, err := store.Load(ctx, initialFunc)
	if err != nil {
		t.Fatalf("Unexpected error from Load(): %v", err)
	}

	if service != mockService {
		t.Error("Expected Load to return the service from initialFunc")
	}
}

func TestOverlayStore_Load_FileDoesNotExist(t *testing.T) {
	_, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// No file created

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.OverlayService {
		return mockService
	}

	service, err := store.Load(ctx, initialFunc)
	if err != nil {
		t.Fatalf("Unexpected error from Load() when file doesn't exist: %v", err)
	}

	if service != mockService {
		t.Error("Expected Load to return the service from initialFunc")
	}
}

func TestOverlayStore_Load_InvalidTOML(t *testing.T) {
	tempDir, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Create an invalid TOML file
	configPath := filepath.Join(tempDir, "overlay.v4.toml")
	err := os.WriteFile(configPath, []byte("invalid toml content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid TOML file: %v", err)
	}

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.OverlayService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() with invalid TOML, got nil")
	}
}

func TestOverlayStore_Load_AddPatternError(t *testing.T) {
	tempDir, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Create test patterns
	pattern := workspace.OverlayPattern{
		Pattern: "*.go",
		Files: []workspace.OverlayFile{
			{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
		},
	}
	patterns := []workspace.OverlayPattern{pattern}

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "overlay.v4.toml")
	createOverlayStoreTestTOMLFile(t, configPath, patterns)

	// Setup mock expectations
	expectedErr := errors.New("test error")
	mockService.EXPECT().AddPattern("*.go", []workspace.OverlayFile{
		{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
	}).Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.OverlayService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err == nil || !errors.Is(errors.Unwrap(err), expectedErr) {
		t.Fatalf("Expected error containing %v from Load(), got %v", expectedErr, err)
	}
}

func TestOverlayStore_Save_NoChanges(t *testing.T) {
	_, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}
}

func TestOverlayStore_Save_WithChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Create test patterns
	patterns := []workspace.OverlayPattern{
		{
			Pattern: "*.go",
			Files: []workspace.OverlayFile{
				{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
			},
		},
	}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(true)
	mockService.EXPECT().GetPatterns().Return(patterns)
	mockService.EXPECT().MarkSaved()

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "overlay.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}

func TestOverlayStore_Save_ForceWithoutChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupOverlayStoreTestEnvironment(t)
	defer cleanup()

	// Create test patterns
	patterns := []workspace.OverlayPattern{
		{
			Pattern: "*.go",
			Files: []workspace.OverlayFile{
				{SourcePath: "/src/file1.go", TargetPath: "/dest/file1.go"},
			},
		},
	}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)
	mockService.EXPECT().GetPatterns().Return(patterns)
	mockService.EXPECT().MarkSaved()

	// Call Save with force=true
	ctx := context.Background()
	err := store.Save(ctx, mockService, true)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "overlay.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}
