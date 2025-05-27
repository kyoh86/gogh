package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	testtarget "github.com/kyoh86/gogh/v4/infra/filesystem"
)

func TestLayoutServiceGetRoot(t *testing.T) {
	// Setup
	root := "/test/root"
	layout := testtarget.NewLayoutService(root)

	// Test
	if layout.GetRoot() != string(root) {
		t.Errorf("Expected root %s, got %s", root, layout.GetRoot())
	}
}

func TestLayoutServiceMatch(t *testing.T) {
	// Setup
	root := "/test/root"
	layout := testtarget.NewLayoutService(root)

	testCases := []struct {
		path          string
		expectSuccess bool
		expectedHost  string
		expectedOwner string
		expectedName  string
	}{
		{
			path:          "/test/root/github.com/kyoh86/gogh",
			expectSuccess: true,
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			path:          "/test/root/github.com/kyoh86/gogh/cmd",
			expectSuccess: true, // Should match with subdirectory
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			path:          "/test/root/github.com/kyoh86", // Not enough components
			expectSuccess: false,
		},
		{
			path:          "/other/path/github.com/kyoh86/gogh", // Not under root
			expectSuccess: false,
		},
	}

	for _, tc := range testCases {
		ref, err := layout.Match(tc.path)

		if tc.expectSuccess {
			if err != nil {
				t.Errorf("Expected success for path %s, got error: %v", tc.path, err)
				continue
			}

			if ref.Host() != tc.expectedHost {
				t.Errorf("For path %s, expected host %s, got %s", tc.path, tc.expectedHost, ref.Host())
			}

			if ref.Owner() != tc.expectedOwner {
				t.Errorf("For path %s, expected owner %s, got %s", tc.path, tc.expectedOwner, ref.Owner())
			}

			if ref.Name() != tc.expectedName {
				t.Errorf("For path %s, expected name %s, got %s", tc.path, tc.expectedName, ref.Name())
			}
		} else {
			if err == nil {
				t.Errorf("Expected error for path %s, got success with ref %v", tc.path, ref)
			}

			if err != workspace.ErrNotMatched {
				t.Errorf("Expected ErrNotMatched for path %s, got %v", tc.path, err)
			}
		}
	}
}

func TestLayoutServiceExactMatch(t *testing.T) {
	// Setup
	root := "/test/root"
	layout := testtarget.NewLayoutService(root)

	testCases := []struct {
		path          string
		expectSuccess bool
		expectedHost  string
		expectedOwner string
		expectedName  string
	}{
		{
			path:          "/test/root/github.com/kyoh86/gogh",
			expectSuccess: true,
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
		},
		{
			path:          "/test/root/github.com/kyoh86/gogh/cmd", // Subdirectory shouldn't match
			expectSuccess: false,
		},
		{
			path:          "/test/root/github.com/kyoh86", // Not enough components
			expectSuccess: false,
		},
		{
			path:          "/other/path/github.com/kyoh86/gogh", // Not under root
			expectSuccess: false,
		},
	}

	for _, tc := range testCases {
		ref, err := layout.ExactMatch(tc.path)

		if tc.expectSuccess {
			if err != nil {
				t.Errorf("Expected success for path %s, got error: %v", tc.path, err)
				continue
			}

			if ref.Host() != tc.expectedHost {
				t.Errorf("For path %s, expected host %s, got %s", tc.path, tc.expectedHost, ref.Host())
			}

			if ref.Owner() != tc.expectedOwner {
				t.Errorf("For path %s, expected owner %s, got %s", tc.path, tc.expectedOwner, ref.Owner())
			}

			if ref.Name() != tc.expectedName {
				t.Errorf("For path %s, expected name %s, got %s", tc.path, tc.expectedName, ref.Name())
			}
		} else {
			if err == nil {
				t.Errorf("Expected error for path %s, got success with ref %v", tc.path, ref)
			}

			if err != workspace.ErrNotMatched {
				t.Errorf("Expected ErrNotMatched for path %s, got %v", tc.path, err)
			}
		}
	}
}

func TestLayoutServicePathFor(t *testing.T) {
	// Setup
	root := "/test/root"
	layout := testtarget.NewLayoutService(root)

	// Create a reference
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	// Test
	expected := filepath.Join("/test/root", "github.com", "kyoh86", "gogh")
	path := layout.PathFor(ref)

	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestLayoutServiceCreateAndDeleteRepository(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "layout-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup
	root := tmpDir
	layout := testtarget.NewLayoutService(root)

	// Create a reference
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	// Test CreateRepositoryFolder
	createdPath, err := layout.CreateRepositoryFolder(ref)
	if err != nil {
		t.Fatalf("Failed to create repository folder: %v", err)
	}

	expected := filepath.Join(tmpDir, "github.com", "kyoh86", "gogh")
	if createdPath != expected {
		t.Errorf("Expected created path %s, got %s", expected, createdPath)
	}

	// Check if directory exists
	info, err := os.Stat(createdPath)
	if err != nil {
		t.Errorf("Failed to stat created directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}

	// Test DeleteRepository
	err = layout.DeleteRepository(ref)
	if err != nil {
		t.Fatalf("Failed to delete repository: %v", err)
	}

	// Check if directory is gone
	_, err = os.Stat(createdPath)
	if !os.IsNotExist(err) {
		t.Errorf("Repository directory still exists after deletion")
	}
}
