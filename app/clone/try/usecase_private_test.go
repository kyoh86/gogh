package try

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"go.uber.org/mock/gomock"
)

func TestValidateExistingRepoStructure(t *testing.T) {
	tests := []struct {
		name             string
		setup            func(ctrl *gomock.Controller) (string, *git_mock.MockGitService)
		requestStructure config.RepositoryStructure
		expectError      bool
		errorContains    string
	}{
		{
			name: "directory does not exist - no conflict",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				return "/nonexistent/path", nil
			},
			requestStructure: config.StructureWorktree,
			expectError:      false,
		},
		{
			name: "normal structure, requesting worktree - error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(false, nil)
				// Create temp directory for test
				tmpDir := t.TempDir()
				return tmpDir, gitService
			},
			requestStructure: config.StructureWorktree,
			expectError:      true,
			errorContains:    "Cannot clone with --structure=worktree flag",
		},
		{
			name: "worktree structure, requesting normal - error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(true, nil)
				// Create temp directory with .worktree subdirectory
				tmpDir := t.TempDir()
				worktreeDir := filepath.Join(tmpDir, ".worktree")
				if err := os.Mkdir(worktreeDir, 0o755); err != nil {
					t.Fatal(err)
				}
				return tmpDir, gitService
			},
			requestStructure: config.StructureNormal,
			expectError:      true,
			errorContains:    "Cannot clone with --structure=normal flag",
		},
		{
			name: "worktree structure, requesting worktree - no error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(true, nil)
				// Create temp directory with .worktree subdirectory
				tmpDir := t.TempDir()
				worktreeDir := filepath.Join(tmpDir, ".worktree")
				if err := os.Mkdir(worktreeDir, 0o755); err != nil {
					t.Fatal(err)
				}
				return tmpDir, gitService
			},
			requestStructure: config.StructureWorktree,
			expectError:      false,
		},
		{
			name: "normal structure, requesting normal - no error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(false, nil)
				// Create temp directory
				tmpDir := t.TempDir()
				return tmpDir, gitService
			},
			requestStructure: config.StructureNormal,
			expectError:      false,
		},
		{
			name: "bare repo without .worktree, requesting normal - no error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(true, nil)
				// Create temp directory without .worktree subdirectory
				tmpDir := t.TempDir()
				return tmpDir, gitService
			},
			requestStructure: config.StructureNormal,
			expectError:      false,
		},
		{
			name: "not a git repository - no error",
			setup: func(ctrl *gomock.Controller) (string, *git_mock.MockGitService) {
				gitService := git_mock.NewMockGitService(ctrl)
				gitService.EXPECT().IsBare(gomock.Any(), gomock.Any()).Return(false, errors.New("not a git repository"))
				// Create temp directory
				tmpDir := t.TempDir()
				return tmpDir, gitService
			},
			requestStructure: config.StructureWorktree,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			localPath, gitService := tt.setup(ctrl)

			err := validateExistingRepoStructure(context.Background(), gitService, localPath, tt.requestStructure.IsWorktree())

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorContains)
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
