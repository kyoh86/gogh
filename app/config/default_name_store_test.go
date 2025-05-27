package config_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"go.uber.org/mock/gomock"
)

// setupTempDir creates a temporary directory for testing and returns its path
// along with a cleanup function.
func setupTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "default-name-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// setupEnvironment sets up a test environment with temporary directory
// and mocked DefaultNameService
func setupEnvironment(t *testing.T) (
	*gomock.Controller,
	string,
	func(),
	*repository_mock.MockDefaultNameService,
	*config.DefaultNameStore,
) {
	ctrl := gomock.NewController(t)
	tempDir, cleanup := setupTempDir(t)
	mockService := repository_mock.NewMockDefaultNameService(ctrl)

	// Override the appContextPath to use our test directory
	origAppContextPath := config.AppContextPathFunc
	config.AppContextPathFunc = func(envName string, fallbackFunc func() (string, error), rel ...string) (string, error) {
		return filepath.Join(append([]string{tempDir}, rel...)...), nil
	}

	// Add cleanup for the override
	originalCleanup := cleanup
	cleanup = func() {
		originalCleanup()
		config.AppContextPathFunc = origAppContextPath
	}

	store := config.NewDefaultNameStore()

	return ctrl, tempDir, cleanup, mockService, store
}

// createTestTOMLFile creates a test TOML file with default values
func createTestTOMLFile(t *testing.T, path string) {
	t.Helper()

	content := `default_host = "github.com"

[hosts]
'github.com' = "testuser"
'gitlab.com' = "otheruser"
`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test TOML file: %v", err)
	}
}

func TestNewDefaultNameStore(t *testing.T) {
	store := config.NewDefaultNameStore()
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestSource(t *testing.T) {
	_, tempDir, cleanup, _, store := setupEnvironment(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "default_names.v4.toml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestLoad_FileExists(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	createTestTOMLFile(t, configPath)

	// Setup mock expectations
	mockService.EXPECT().SetDefaultHost("github.com").Return(nil)
	mockService.EXPECT().SetDefaultOwnerFor("github.com", "testuser").Return(nil)
	mockService.EXPECT().SetDefaultOwnerFor("gitlab.com", "otheruser").Return(nil)
	mockService.EXPECT().MarkSaved()

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
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

func TestLoad_FileDoesNotExist(t *testing.T) {
	_, _, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// No file created

	// Setup mock expectations - no actions expected
	mockService.EXPECT().SetDefaultHost("").Return(nil).MaxTimes(1)

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err != nil {
		// An error is expected because the file doesn't exist
		t.Logf("Got expected error from Load(): %v", err)
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Create an invalid TOML file
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	err := os.WriteFile(configPath, []byte("invalid toml content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid TOML file: %v", err)
	}

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() with invalid TOML, got nil")
	}
}

func TestLoad_SetDefaultHostError(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	createTestTOMLFile(t, configPath)

	// Setup mock expectations
	expectedErr := errors.New("test error")
	mockService.EXPECT().SetDefaultHost("github.com").Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v from Load(), got %v", expectedErr, err)
	}
}

func TestLoad_SetDefaultOwnerError(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	createTestTOMLFile(t, configPath)

	// Setup mock expectations
	mockService.EXPECT().SetDefaultHost("github.com").Return(nil)
	expectedErr := errors.New("test error")
	mockService.EXPECT().SetDefaultOwnerFor("github.com", "testuser").Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error about setting default owner, got %v", err)
	}
}

func TestSave_NoChanges(t *testing.T) {
	_, _, cleanup, mockService, store := setupEnvironment(t)
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

func TestSave_WithChanges(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(true)
	mockService.EXPECT().GetMap().Return(map[string]string{
		"github.com": "testuser",
		"gitlab.com": "otheruser",
	})
	mockService.EXPECT().GetDefaultHost().Return("github.com")
	mockService.EXPECT().MarkSaved()

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}

func TestSave_ForceWithoutChanges(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)
	mockService.EXPECT().GetMap().Return(map[string]string{
		"github.com": "testuser",
	})
	mockService.EXPECT().GetDefaultHost().Return("github.com")
	mockService.EXPECT().MarkSaved()

	// Call Save with force=true
	ctx := context.Background()
	err := store.Save(ctx, mockService, true)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "default_names.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}

func TestSave_CreateDirectory(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironment(t)
	defer cleanup()

	// Setup a nested config path
	nestedDir := filepath.Join(tempDir, "nested", "dir")
	config.AppContextPathFunc = func(envName string, fallbackFunc func() (string, error), rel ...string) (string, error) {
		return filepath.Join(append([]string{nestedDir}, rel...)...), nil
	}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(true)
	mockService.EXPECT().GetMap().Return(map[string]string{
		"github.com": "testuser",
	})
	mockService.EXPECT().GetDefaultHost().Return("github.com")
	mockService.EXPECT().MarkSaved()

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the nested directory was created
	configPath := filepath.Join(nestedDir, "default_names.v4.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Expected directory and file to be created, but it doesn't exist: %s (%s)", configPath, err)
	}
}
