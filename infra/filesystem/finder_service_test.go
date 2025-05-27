package filesystem_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"github.com/kyoh86/gogh/v4/infra/filesystem"
	"go.uber.org/mock/gomock"
)

// setupTestEnvironment creates a test environment with temporary directories
func setupTestEnvironment(t *testing.T) string {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "finder-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

func TestFindByReference(t *testing.T) {
	// Set up test environment
	tmpDir := setupTestEnvironment(t)
	defer os.RemoveAll(tmpDir)

	// Set up test roots
	root1 := filepath.Join(tmpDir, "root1")
	root2 := filepath.Join(tmpDir, "root2")

	// Create test repository directories
	repoPath1 := filepath.Join(string(root1), "github.com", "kyoh86", "gogh")
	repoPath2 := filepath.Join(string(root2), "gitlab.com", "user", "project")

	if err := os.MkdirAll(repoPath1, 0755); err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}
	if err := os.MkdirAll(repoPath2, 0755); err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name          string
		ref           repository.Reference
		expectFound   bool
		expectedHost  string
		expectedOwner string
		expectedName  string
	}{
		{
			name:          "existing repository in first root",
			ref:           repository.NewReference("github.com", "kyoh86", "gogh"),
			expectFound:   true,
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			name:          "existing repository in second root",
			ref:           repository.NewReference("gitlab.com", "user", "project"),
			expectFound:   true,
			expectedHost:  "gitlab.com",
			expectedOwner: "user",
			expectedName:  "project",
		},
		{
			name:        "non-existing repository",
			ref:         repository.NewReference("github.com", "nonexist", "repo"),
			expectFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new controller and mock for each test case
			ctrl := gomock.NewController(t)
			mockWS := workspace_mock.NewMockWorkspaceService(ctrl)

			// Create layout services
			mockLayout1 := filesystem.NewLayoutService(root1)
			mockLayout2 := filesystem.NewLayoutService(root2)

			// Create a FinderService
			finder := filesystem.NewFinderService()

			// Set up mock expectations with the new mock
			mockWS.EXPECT().GetRoots().Return([]workspace.Root{root1, root2}).AnyTimes()
			mockWS.EXPECT().GetLayoutFor(root1).Return(mockLayout1).AnyTimes()
			mockWS.EXPECT().GetLayoutFor(root2).Return(mockLayout2).AnyTimes()

			loc, err := finder.FindByReference(ctx, mockWS, tc.ref)

			if tc.expectFound {
				if err != nil {
					t.Errorf("Expected to find repository, got error: %v", err)
					return
				}

				if loc == nil {
					t.Error("Expected non-nil location, got nil")
					return
				}

				if loc.Host() != tc.expectedHost {
					t.Errorf("Expected host %s, got %s", tc.expectedHost, loc.Host())
				}

				if loc.Owner() != tc.expectedOwner {
					t.Errorf("Expected owner %s, got %s", tc.expectedOwner, loc.Owner())
				}

				if loc.Name() != tc.expectedName {
					t.Errorf("Expected name %s, got %s", tc.expectedName, loc.Name())
				}
			} else {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				if !errors.Is(err, workspace.ErrNotMatched) {
					t.Errorf("Expected ErrNotMatched, got %v", err)
				}
			}
		})
	}
}

func TestFindByPath(t *testing.T) {
	// Set up test environment
	tmpDir := setupTestEnvironment(t)
	defer os.RemoveAll(tmpDir)

	// Set up test root
	root1 := filepath.Join(tmpDir, "root1")

	// Create test repository directory
	repoPath := filepath.Join(string(root1), "github.com", "kyoh86", "gogh")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name          string
		path          string
		setupMock     func(mockWS *workspace_mock.MockWorkspaceService)
		expectFound   bool
		expectedHost  string
		expectedOwner string
		expectedName  string
	}{
		{
			name: "absolute path",
			path: repoPath,
			setupMock: func(mockWS *workspace_mock.MockWorkspaceService) {
				mockLayout := filesystem.NewLayoutService(root1)
				mockWS.EXPECT().GetRoots().Return([]workspace.Root{root1}).AnyTimes()
				mockWS.EXPECT().GetLayoutFor(root1).Return(mockLayout).AnyTimes()
				mockWS.EXPECT().GetPrimaryRoot().Return(root1).AnyTimes()
				mockWS.EXPECT().GetPrimaryLayout().Return(mockLayout).AnyTimes()
			},
			expectFound:   true,
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			name: "relative path",
			path: "github.com/kyoh86/gogh",
			setupMock: func(mockWS *workspace_mock.MockWorkspaceService) {
				mockLayout := filesystem.NewLayoutService(root1)
				mockWS.EXPECT().GetRoots().Return([]workspace.Root{root1}).AnyTimes()
				mockWS.EXPECT().GetLayoutFor(root1).Return(mockLayout).AnyTimes()
			},
			expectFound:   true,
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			name: "non-existing path",
			path: "github.com/nonexist/repo",
			setupMock: func(mockWS *workspace_mock.MockWorkspaceService) {
				mockLayout := filesystem.NewLayoutService(root1)
				mockWS.EXPECT().GetRoots().Return([]workspace.Root{root1}).AnyTimes()
				mockWS.EXPECT().GetLayoutFor(root1).Return(mockLayout).AnyTimes()
			},
			expectFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new controller and mock for each test case
			ctrl := gomock.NewController(t)
			mockWS := workspace_mock.NewMockWorkspaceService(ctrl)
			tc.setupMock(mockWS)

			// Create a FinderService
			finder := filesystem.NewFinderService()

			loc, err := finder.FindByPath(ctx, mockWS, tc.path)

			if tc.expectFound {
				if err != nil {
					t.Errorf("Expected to find repository, got error: %v", err)
					return
				}

				if loc == nil {
					t.Error("Expected non-nil location, got nil")
					return
				}

				if loc.Host() != tc.expectedHost {
					t.Errorf("Expected host %s, got %s", tc.expectedHost, loc.Host())
				}

				if loc.Owner() != tc.expectedOwner {
					t.Errorf("Expected owner %s, got %s", tc.expectedOwner, loc.Owner())
				}

				if loc.Name() != tc.expectedName {
					t.Errorf("Expected name %s, got %s", tc.expectedName, loc.Name())
				}
			} else {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				if !errors.Is(err, workspace.ErrNotMatched) {
					t.Errorf("Expected ErrNotMatched, got %v", err)
				}
			}
		})
	}
}

func TestListAllRepository(t *testing.T) {
	// Set up test environment
	tmpDir := setupTestEnvironment(t)
	defer os.RemoveAll(tmpDir)

	// Set up test roots
	root1 := filepath.Join(tmpDir, "root1")
	root2 := filepath.Join(tmpDir, "root2")

	// Create test repository directories
	if err := os.MkdirAll(filepath.Join(string(root1), "github.com", "kyoh86", "gogh"), 0755); err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(string(root2), "gitlab.com", "user", "project"), 0755); err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}

	// Create a new controller and mock
	ctrl := gomock.NewController(t)
	mockWS := workspace_mock.NewMockWorkspaceService(ctrl)

	// Create layout services
	mockLayout1 := filesystem.NewLayoutService(root1)
	mockLayout2 := filesystem.NewLayoutService(root2)

	// Create a FinderService
	finder := filesystem.NewFinderService()

	// Create context
	ctx := context.Background()

	// Set up mock expectations
	mockWS.EXPECT().GetRoots().Return([]workspace.Root{root1, root2}).AnyTimes()
	mockWS.EXPECT().GetLayoutFor(root1).Return(mockLayout1).AnyTimes()
	mockWS.EXPECT().GetLayoutFor(root2).Return(mockLayout2).AnyTimes()

	// List repositories
	var found []string
	var errs []error

	opts := workspace.ListOptions{}
	for loc, err := range finder.ListAllRepository(ctx, mockWS, opts) {
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if loc != nil {
			found = append(found, loc.Path())
		}
	}

	// Verify no errors occurred
	if len(errs) > 0 {
		t.Errorf("Unexpected errors: %v", errs)
	}

	// Verify that both repositories were found
	if len(found) != 2 {
		t.Errorf("Expected to find 2 repositories, got %d: %v", len(found), found)
	}
}
