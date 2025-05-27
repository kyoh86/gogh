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
	"go.uber.org/mock/gomock"
)

// setupWorkspaceStoreTestTempDir creates a temporary directory for testing and returns its path
// along with a cleanup function.
func setupWorkspaceStoreTestTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "workspace-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// setupWorkspaceStoreTestEnvironment sets up a test environment with temporary directory
// and mocked WorkspaceService
func setupWorkspaceStoreTestEnvironment(t *testing.T) (
	string,
	func(),
	*workspace_mock.MockWorkspaceService,
	*config.WorkspaceStore,
) {
	ctrl := gomock.NewController(t)
	tempDir, cleanup := setupWorkspaceStoreTestTempDir(t)
	mockService := workspace_mock.NewMockWorkspaceService(ctrl)

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

	store := config.NewWorkspaceStore()

	return tempDir, cleanup, mockService, store
}

// createWorkspaceStoreTestTOMLFile creates a test TOML file with workspace settings
func createWorkspaceStoreTestTOMLFile(t *testing.T, path string, roots []workspace.Root, primaryRoot workspace.Root) {
	t.Helper()

	content := "roots = [\n"
	for _, root := range roots {
		content += `  "` + root + `",` + "\n"
	}
	content += "]\n"
	content += `primary_root = "` + primaryRoot + `"`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test TOML file: %v", err)
	}
}

func TestNewWorkspaceStore(t *testing.T) {
	store := config.NewWorkspaceStore()
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestWorkspaceStore_Source(t *testing.T) {
	tempDir, cleanup, _, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "workspace.v4.toml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestWorkspaceStore_Load_FileExists(t *testing.T) {
	tempDir, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// Create test roots
	root1 := filepath.Join(tempDir, "root1")
	root2 := filepath.Join(tempDir, "root2")
	roots := []workspace.Root{root1, root2}

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "workspace.v4.toml")
	createWorkspaceStoreTestTOMLFile(t, configPath, roots, root1)

	// Setup mock expectations
	mockService.EXPECT().AddRoot(root1, true).Return(nil)
	mockService.EXPECT().AddRoot(root2, false).Return(nil)
	mockService.EXPECT().MarkSaved()

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.WorkspaceService {
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

func TestWorkspaceStore_Load_FileDoesNotExist(t *testing.T) {
	_, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// No file created

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.WorkspaceService {
		return mockService
	}

	if _, err := store.Load(ctx, initialFunc); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Expected error from Load() as file doesn't exist: %v", err)
	}
}

func TestWorkspaceStore_Load_InvalidTOML(t *testing.T) {
	tempDir, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// Create an invalid TOML file
	configPath := filepath.Join(tempDir, "workspace.v4.toml")
	err := os.WriteFile(configPath, []byte("invalid toml content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid TOML file: %v", err)
	}

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.WorkspaceService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() with invalid TOML, got nil")
	}
}

func TestWorkspaceStore_Load_AddRootError(t *testing.T) {
	tempDir, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// Create test roots
	root1 := filepath.Join(tempDir, "root1")
	root2 := filepath.Join(tempDir, "root2")
	roots := []workspace.Root{root1, root2}

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "workspace.v4.toml")
	createWorkspaceStoreTestTOMLFile(t, configPath, roots, root1)

	// Setup mock expectations
	expectedErr := errors.New("test error")
	mockService.EXPECT().AddRoot(root1, true).Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() workspace.WorkspaceService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err != expectedErr {
		t.Fatalf("Expected error %v from Load(), got %v", expectedErr, err)
	}
}

func TestWorkspaceStore_Save_NoChanges(t *testing.T) {
	_, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
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

func TestWorkspaceStore_Save_WithChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// Create test roots
	root1 := filepath.Join(tempDir, "root1")
	root2 := filepath.Join(tempDir, "root2")
	roots := []workspace.Root{root1, root2}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(true)
	mockService.EXPECT().GetRoots().Return(roots)
	mockService.EXPECT().GetPrimaryRoot().Return(root1)
	mockService.EXPECT().MarkSaved()

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "workspace.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}

func TestWorkspaceStore_Save_ForceWithoutChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupWorkspaceStoreTestEnvironment(t)
	defer cleanup()

	// Create test roots
	root1 := filepath.Join(tempDir, "root1")
	roots := []workspace.Root{root1}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)
	mockService.EXPECT().GetRoots().Return(roots)
	mockService.EXPECT().GetPrimaryRoot().Return(root1)
	mockService.EXPECT().MarkSaved()

	// Call Save with force=true
	ctx := context.Background()
	err := store.Save(ctx, mockService, true)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "workspace.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}
