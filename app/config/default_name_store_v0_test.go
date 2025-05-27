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

// setupTempDirV0 creates a temporary directory for testing and returns its path
// along with a cleanup function.
func setupTempDirV0(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "default-name-store-v0-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// setupEnvironmentV0 sets up a test environment with temporary directory
// and mocked DefaultNameService for V0 store
func setupEnvironmentV0(t *testing.T) (
	*gomock.Controller,
	string,
	func(),
	*repository_mock.MockDefaultNameService,
	*config.DefaultNameStoreV0,
) {
	ctrl := gomock.NewController(t)
	tempDir, cleanup := setupTempDirV0(t)
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

	store := config.NewDefaultNameStoreV0()

	return ctrl, tempDir, cleanup, mockService, store
}

// createTestYAMLFile creates a test YAML file with default values
func createTestYAMLFile(t *testing.T, path string) {
	t.Helper()

	content := `default_host: github.com
hosts:
  github.com:
    default_owner: testuser
  gitlab.com:
    default_owner: otheruser
`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test YAML file: %v", err)
	}
}

func TestNewDefaultNameStoreV0(t *testing.T) {
	store := config.NewDefaultNameStoreV0()
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestSourceV0(t *testing.T) {
	_, tempDir, cleanup, _, store := setupEnvironmentV0(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "tokens.yaml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestLoadV0_FileExists(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create a test YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	createTestYAMLFile(t, configPath)

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

func TestLoadV0_FileDoesNotExist(t *testing.T) {
	_, _, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// No file created

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() when file doesn't exist, got nil")
	}
}

func TestLoadV0_InvalidYAML(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create an invalid YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: - - -"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid YAML file: %v", err)
	}

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() with invalid YAML, got nil")
	}
}

func TestLoadV0_SetDefaultHostError(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create a test YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	createTestYAMLFile(t, configPath)

	// Setup mock expectations
	expectedErr := errors.New("test error")
	mockService.EXPECT().SetDefaultHost("github.com").Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() repository.DefaultNameService {
		return mockService
	}

	if _, err := store.Load(ctx, initialFunc); !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v from Load(), got %v", expectedErr, err)
	}
}

func TestLoadV0_SetDefaultOwnerError(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create a test YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	createTestYAMLFile(t, configPath)

	// Setup mock expectations
	mockService.EXPECT().SetDefaultHost("github.com").Return(nil)
	expectedErr := errors.New("test error")
	mockService.EXPECT().SetDefaultOwnerFor("github.com", "testuser").Return(expectedErr).AnyTimes()
	mockService.EXPECT().SetDefaultOwnerFor("gitlab.com", "otheruser").Return(expectedErr).AnyTimes()

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

func TestLoadV0_EmptyOwner(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create a YAML file with an entry that has an empty default_owner
	configPath := filepath.Join(tempDir, "tokens.yaml")
	content := `default_host: github.com
hosts:
  github.com:
    default_owner: testuser
  gitlab.com:
    default_owner: ""
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	// Setup mock expectations
	mockService.EXPECT().SetDefaultHost("github.com").Return(nil)
	mockService.EXPECT().SetDefaultOwnerFor("github.com", "testuser").Return(nil)
	// Note: gitlab.com entry should be skipped due to empty owner
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

func TestLoadV0_ComplexYAML(t *testing.T) {
	_, tempDir, cleanup, mockService, store := setupEnvironmentV0(t)
	defer cleanup()

	// Create a more complex YAML file with additional fields
	configPath := filepath.Join(tempDir, "tokens.yaml")
	content := `default_host: github.com
hosts:
  github.com:
    default_owner: testuser
    extra_field: should be ignored
  gitlab.com:
    default_owner: otheruser
    another_field: also ignored
extra_section:
  some_value: this should be ignored too
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write complex YAML file: %v", err)
	}

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
