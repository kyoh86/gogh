package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/infra/filesystem"
)

func TestLayoutService_GetRoot(t *testing.T) {
	root := "/home/user/repos"
	layout := filesystem.NewLayoutService(root)

	if got := layout.GetRoot(); got != root {
		t.Errorf("GetRoot() = %q, want %q", got, root)
	}
}

func TestLayoutService_Match(t *testing.T) {
	root := "/home/user/repos"
	layout := filesystem.NewLayoutService(root)

	tests := []struct {
		name    string
		path    string
		want    *repository.Reference
		wantErr bool
	}{
		{
			name: "exact repository path",
			path: "/home/user/repos/github.com/kyoh86/gogh",
			want: func() *repository.Reference { r := repository.NewReference("github.com", "kyoh86", "gogh"); return &r }(),
		},
		{
			name: "nested path in repository",
			path: "/home/user/repos/github.com/kyoh86/gogh/cmd/main.go",
			want: func() *repository.Reference { r := repository.NewReference("github.com", "kyoh86", "gogh"); return &r }(),
		},
		{
			name:    "path outside root",
			path:    "/home/other/repos/github.com/kyoh86/gogh",
			wantErr: true,
		},
		{
			name:    "path with too few components",
			path:    "/home/user/repos/github.com",
			wantErr: true,
		},
		{
			name:    "relative path with ..",
			path:    "/home/user/repos/../other/github.com/kyoh86/gogh",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := layout.Match(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("Match() returned nil without error")
				return
			}
			if !tt.wantErr && got.String() != tt.want.String() {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutService_ExactMatch(t *testing.T) {
	root := "/home/user/repos"
	layout := filesystem.NewLayoutService(root)

	tests := []struct {
		name    string
		path    string
		want    *repository.Reference
		wantErr bool
	}{
		{
			name: "exact repository path",
			path: "/home/user/repos/github.com/kyoh86/gogh",
			want: func() *repository.Reference { r := repository.NewReference("github.com", "kyoh86", "gogh"); return &r }(),
		},
		{
			name:    "nested path in repository",
			path:    "/home/user/repos/github.com/kyoh86/gogh/cmd/main.go",
			wantErr: true,
		},
		{
			name:    "path with too few components",
			path:    "/home/user/repos/github.com/kyoh86",
			wantErr: true,
		},
		{
			name:    "path with too many components",
			path:    "/home/user/repos/github.com/kyoh86/gogh/extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := layout.ExactMatch(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExactMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ExactMatch() returned nil without error")
				return
			}
			if !tt.wantErr && got.String() != tt.want.String() {
				t.Errorf("ExactMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutService_PathFor(t *testing.T) {
	root := "/home/user/repos"
	layout := filesystem.NewLayoutService(root)

	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	want := filepath.Join(root, "github.com", "kyoh86", "gogh")

	if got := layout.PathFor(ref); got != want {
		t.Errorf("PathFor() = %q, want %q", got, want)
	}
}

func TestLayoutService_CreateRepositoryFolder(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "layout_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	layout := filesystem.NewLayoutService(tmpDir)
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	// Create repository folder
	path, err := layout.CreateRepositoryFolder(ref)
	if err != nil {
		t.Fatalf("CreateRepositoryFolder() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "github.com", "kyoh86", "gogh")
	if path != expectedPath {
		t.Errorf("CreateRepositoryFolder() path = %q, want %q", path, expectedPath)
	}

	// Check if directory was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("CreateRepositoryFolder() did not create directory")
	}
}

func TestLayoutService_DeleteRepository(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "layout_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	layout := filesystem.NewLayoutService(tmpDir)
	ref := repository.NewReference("github.com", "kyoh86", "gogh")

	// Create repository folder first
	path, err := layout.CreateRepositoryFolder(ref)
	if err != nil {
		t.Fatalf("CreateRepositoryFolder() error = %v", err)
	}

	// Create a file in the repository
	testFile := filepath.Join(path, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Delete repository
	if err := layout.DeleteRepository(ref); err != nil {
		t.Fatalf("DeleteRepository() error = %v", err)
	}

	// Check if directory was deleted
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("DeleteRepository() did not delete directory")
	}
}

func TestLayoutService_Error(t *testing.T) {
	t.Run("ErrNotMatched", func(t *testing.T) {
		err := workspace.ErrNotMatched
		if err.Error() != "repository not matched for a layout" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
