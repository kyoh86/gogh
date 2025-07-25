package config_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestWorkspaceStoreV0_Load(t *testing.T) {
	// Setup temporary directory for test files
	tempDir, err := os.MkdirTemp("", "gogh-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock AppContextPathFunc to return our test file path
	origAppContextPathFunc := testtarget.AppContextPathFunc
	defer func() { testtarget.AppContextPathFunc = origAppContextPathFunc }()

	configPath := filepath.Join(tempDir, "config.yaml")
	testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return configPath, nil
	}

	t.Run("valid config", func(t *testing.T) {
		// Create valid YAML file
		validYAML := `
roots:
  - /path/to/root1
  - /path/to/root2
`
		err := os.WriteFile(configPath, []byte(validYAML), 0o644)
		if err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		// Create mock service using GoMock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := workspace_mock.NewMockWorkspaceService(ctrl)

		// Setup expectations
		mockSvc.EXPECT().AddRoot(workspace.Root("/path/to/root1"), true).Return(nil)
		mockSvc.EXPECT().AddRoot(workspace.Root("/path/to/root2"), false).Return(nil)
		mockSvc.EXPECT().MarkSaved()

		// Create store and load
		store := testtarget.NewWorkspaceStoreV0()
		result, err := store.Load(context.Background(), func() workspace.WorkspaceService {
			return mockSvc
		})
		if err != nil {
			t.Fatalf("Load failed with error: %v", err)
		}
		if result != mockSvc {
			t.Errorf("expected result to be the mock service")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		// Remove config file
		os.Remove(configPath)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := workspace_mock.NewMockWorkspaceService(ctrl)

		store := testtarget.NewWorkspaceStoreV0()
		_, err := store.Load(context.Background(), func() workspace.WorkspaceService {
			return mockSvc
		})

		if err == nil {
			t.Errorf("expected error when file not found, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		// Create invalid YAML
		invalidYAML := `
roots:
  - /path/to/root1
  - 
  invalid: content
`
		err := os.WriteFile(configPath, []byte(invalidYAML), 0o644)
		if err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSvc := workspace_mock.NewMockWorkspaceService(ctrl)

		store := testtarget.NewWorkspaceStoreV0()
		_, err = store.Load(context.Background(), func() workspace.WorkspaceService {
			return mockSvc
		})

		if err == nil {
			t.Errorf("expected error for invalid YAML, got nil")
		}
	})
}

func TestWorkspaceStoreV0_Source(t *testing.T) {
	origAppContextPathFunc := testtarget.AppContextPathFunc
	defer func() { testtarget.AppContextPathFunc = origAppContextPathFunc }()

	expectedPath := "/expected/path/config.yaml"
	testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		if envar != "GOGH_CONFIG_PATH" {
			t.Errorf("expected envar to be 'GOGH_CONFIG_PATH', got %q", envar)
		}
		if !reflect.DeepEqual(rel, []string{"config.yaml"}) {
			t.Errorf("expected rel to be []string{\"config.yaml\"}, got %v", rel)
		}
		return expectedPath, nil
	}

	store := testtarget.NewWorkspaceStoreV0()
	path, err := store.Source()
	if err != nil {
		t.Fatalf("Source failed with error: %v", err)
	}
	if path != expectedPath {
		t.Errorf("expected %q, got %q", expectedPath, path)
	}
}
