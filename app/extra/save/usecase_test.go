package save_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/extra/save"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		repoStr   string
		setupMock func(*testing.T, *gomock.Controller) (
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*git_mock.MockGitService,
			*overlay_mock.MockOverlayService,
			*script_mock.MockScriptService,
			*hook_mock.MockHookService,
			*extra_mock.MockExtraService,
			*repository_mock.MockReferenceParser,
		)
		wantErr bool
	}{
		{
			name:    "Successfully save excluded files as auto extra",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				// Create temp directory and files
				tempDir := t.TempDir()

				// Create test files
				gitignoreContent := []byte("*.log\nnode_modules/\n")
				if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), gitignoreContent, 0o644); err != nil {
					t.Fatalf("Failed to create .gitignore: %v", err)
				}
				configContent := []byte("[settings]\nkey = \"value\"\n")
				if err := os.WriteFile(filepath.Join(tempDir, "config.toml"), configContent, 0o644); err != nil {
					t.Fatalf("Failed to create config.toml: %v", err)
				}

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					tempDir,
					"github.com",
					"owner",
					"repo",
				)

				// Parse reference - called twice: once in GetExcludedFiles, once in SaveFiles
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil).Times(2)

				// Find repository - called twice: once in GetExcludedFiles, once in SaveFiles
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil).Times(2)

				// Check existing auto extra
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// List excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, ".gitignore"), nil)
						yield(filepath.Join(tempDir, "config.toml"), nil)
					},
				)

				// Create overlays and hooks
				overlayID1 := uuid.New()
				overlayID2 := uuid.New()
				hookID1 := uuid.New().String()
				hookID2 := uuid.New().String()

				// For .gitignore
				overlayService.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.Name != "Auto extra: .gitignore" {
							t.Errorf("Expected overlay name 'Auto extra: .gitignore', got %s", entry.Name)
						}
						if entry.RelativePath != ".gitignore" {
							t.Errorf("Expected relative path '.gitignore', got %s", entry.RelativePath)
						}
						return overlayID1.String(), nil
					},
				)
				// Mock overlay resolution for hook operation ID
				mockOverlay1 := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay1.EXPECT().UUID().Return(overlayID1)
				overlayService.EXPECT().Get(ctx, overlayID1.String()).Return(mockOverlay1, nil)
				hs.EXPECT().Add(ctx, gomock.Any()).Return(hookID1, nil)

				// For config.toml
				overlayService.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.Name != "Auto extra: config.toml" {
							t.Errorf("Expected overlay name 'Auto extra: config.toml', got %s", entry.Name)
						}
						if entry.RelativePath != "config.toml" {
							t.Errorf("Expected relative path 'config.toml', got %s", entry.RelativePath)
						}
						return overlayID2.String(), nil
					},
				)
				// Mock overlay resolution for hook operation ID
				mockOverlay2 := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay2.EXPECT().UUID().Return(overlayID2)
				overlayService.EXPECT().Get(ctx, overlayID2.String()).Return(mockOverlay2, nil)
				hs.EXPECT().Add(ctx, gomock.Any()).Return(hookID2, nil)

				// Save auto extra
				es.EXPECT().AddAutoExtra(ctx, ref, ref, gomock.Any()).DoAndReturn(
					func(ctx context.Context, ref repository.Reference, targetRef repository.Reference, items []extra.Item) (string, error) {
						if len(items) != 2 {
							t.Errorf("Expected 2 items, got %d", len(items))
						}
						return uuid.New().String(), nil
					},
				)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: false,
		},
		{
			name:    "Invalid repository reference",
			repoStr: "invalid-ref",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				rp.EXPECT().Parse("invalid-ref").Return(nil, errors.New("invalid reference"))

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Repository not found",
			repoStr: "github.com/owner/notfound",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "notfound")

				rp.EXPECT().Parse("github.com/owner/notfound").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(nil, errors.New("not found"))

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Repository found but location is nil",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(nil, nil)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Auto extra already exists",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				tempDir := t.TempDir()

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					tempDir,
					"github.com",
					"owner",
					"repo",
				)

				// Parse reference - called twice: once in GetExcludedFiles, once in SaveFiles
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil).Times(2)
				// Find repository - called twice: once in GetExcludedFiles, once in SaveFiles
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil).Times(2)

				// List excluded files - has files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, "test.txt"), nil)
					},
				)

				// Existing auto extra
				existingExtra := extra.NewAutoExtra(
					uuid.New().String(),
					ref,
					ref,
					[]extra.Item{},
					time.Now(),
				)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(existingExtra, nil)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "No excluded files in repository",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					"/path/to/repo",
					"github.com",
					"owner",
					"repo",
				)

				// Parse reference - called once in GetExcludedFiles only
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				// Find repository - called once in GetExcludedFiles only
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

				// Empty excluded files list
				gs.EXPECT().ListExcludedFiles(ctx, "/path/to/repo", nil).Return(
					func(yield func(string, error) bool) {
						// No files
					},
				)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Overlay creation fails with rollback",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				// Create temp directory and files
				tempDir := t.TempDir()

				// Create test files
				if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte("*.log\n"), 0o644); err != nil {
					t.Fatalf("Failed to create .gitignore: %v", err)
				}
				if err := os.WriteFile(filepath.Join(tempDir, "config.toml"), []byte("[test]\n"), 0o644); err != nil {
					t.Fatalf("Failed to create config.toml: %v", err)
				}

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					tempDir,
					"github.com",
					"owner",
					"repo",
				)

				// Parse reference - called twice: once in GetExcludedFiles, once in SaveFiles
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil).Times(2)
				// Find repository - called twice: once in GetExcludedFiles, once in SaveFiles
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil).Times(2)

				// List excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, ".gitignore"), nil)
						yield(filepath.Join(tempDir, "config.toml"), nil)
					},
				)

				// Check existing auto extra
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// First overlay succeeds
				overlayID1 := uuid.New()
				overlayService.EXPECT().Add(ctx, gomock.Any()).Return(overlayID1.String(), nil)
				// Mock overlay resolution for hook operation ID
				mockOverlay1 := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay1.EXPECT().UUID().Return(overlayID1)
				overlayService.EXPECT().Get(ctx, overlayID1.String()).Return(mockOverlay1, nil)
				hookID1 := uuid.New().String()
				hs.EXPECT().Add(ctx, gomock.Any()).Return(hookID1, nil)

				// Second overlay fails
				overlayService.EXPECT().Add(ctx, gomock.Any()).Return("", errors.New("overlay creation failed"))

				// Rollback expectations
				overlayService.EXPECT().Remove(ctx, overlayID1.String()).Return(nil)
				hs.EXPECT().Remove(ctx, hookID1).Return(nil)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Hook creation fails with rollback",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				// Create temp directory and files
				tempDir := t.TempDir()

				// Create test files
				if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte("*.log\n"), 0o644); err != nil {
					t.Fatalf("Failed to create .gitignore: %v", err)
				}

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					tempDir,
					"github.com",
					"owner",
					"repo",
				)

				// Parse reference - called twice: once in GetExcludedFiles, once in SaveFiles
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil).Times(2)
				// Find repository - called twice: once in GetExcludedFiles, once in SaveFiles
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil).Times(2)

				// List excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, ".gitignore"), nil)
					},
				)

				// Check existing auto extra
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// Overlay succeeds
				overlayID1 := uuid.New()
				overlayService.EXPECT().Add(ctx, gomock.Any()).Return(overlayID1.String(), nil)
				// Mock overlay resolution for hook operation ID
				mockOverlay1 := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay1.EXPECT().UUID().Return(overlayID1)
				overlayService.EXPECT().Get(ctx, overlayID1.String()).Return(mockOverlay1, nil)

				// Hook fails
				hs.EXPECT().Add(ctx, gomock.Any()).Return("", errors.New("hook creation failed"))

				// Rollback expectations
				overlayService.EXPECT().Remove(ctx, overlayID1.String()).Return(nil)

				ss := script_mock.NewMockScriptService(ctrl)
				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ws, fs, gs, overlayService, ss, hs, es, rp := tc.setupMock(t, ctrl)
			uc := testtarget.NewUsecase(ws, fs, gs, overlayService, ss, hs, es, rp)

			err := uc.Execute(ctx, tc.repoStr)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsecase_GetExcludedFiles(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		repoStr   string
		setupMock func(*testing.T, *gomock.Controller) (
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*git_mock.MockGitService,
			*repository_mock.MockReferenceParser,
		)
		want    *testtarget.ExcludedFilesResult
		wantErr bool
	}{
		{
			name:    "Successfully get excluded files",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				tempDir := t.TempDir()
				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(tempDir, "github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

				// List excluded files
				file1 := filepath.Join(tempDir, ".gitignore")
				file2 := filepath.Join(tempDir, "config.local")
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(file1, nil)
						yield(file2, nil)
					},
				)

				return ws, fs, gs, rp
			},
			want:    nil, // Will be set in the test execution
			wantErr: false,
		},
		{
			name:    "Invalid repository reference",
			repoStr: "invalid-repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				rp.EXPECT().Parse("invalid-repo").Return(nil, errors.New("invalid reference"))

				return ws, fs, gs, rp
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Repository not found",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				ref := repository.NewReference("github.com", "owner", "repo")
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(nil, errors.New("not found"))

				return ws, fs, gs, rp
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No excluded files",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				tempDir := t.TempDir()
				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(tempDir, "github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

				// No excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						// Don't yield any files
					},
				)

				return ws, fs, gs, rp
			},
			want:    nil, // Will be set in the test execution
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ws, fs, gs, rp := tc.setupMock(t, ctrl)

			// Create dummy services for the full constructor
			overlayService := overlay_mock.NewMockOverlayService(ctrl)
			ss := script_mock.NewMockScriptService(ctrl)
			hs := hook_mock.NewMockHookService(ctrl)
			es := extra_mock.NewMockExtraService(ctrl)

			uc := testtarget.NewUsecase(ws, fs, gs, overlayService, ss, hs, es, rp)

			got, err := uc.GetExcludedFiles(ctx, tc.repoStr)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetExcludedFiles() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// For successful cases, check the result
			if !tc.wantErr && got != nil {
				// For "Successfully get excluded files" case, check we got 2 files
				if tc.name == "Successfully get excluded files" && len(got.Files) != 2 {
					t.Errorf("GetExcludedFiles() got %d files, want 2 files", len(got.Files))
				}
				// For "No excluded files" case, check we got 0 files
				if tc.name == "No excluded files" && len(got.Files) != 0 {
					t.Errorf("GetExcludedFiles() got %d files, want 0 files", len(got.Files))
				}
				// Check repository path is not empty
				if got.RepositoryPath == "" {
					t.Errorf("GetExcludedFiles() RepositoryPath is empty")
				}
			}
		})
	}
}

func TestUsecase_SaveFiles(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		repoStr   string
		files     []string
		setupMock func(*testing.T, *gomock.Controller, []string) (
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*git_mock.MockGitService,
			*overlay_mock.MockOverlayService,
			*script_mock.MockScriptService,
			*hook_mock.MockHookService,
			*extra_mock.MockExtraService,
			*repository_mock.MockReferenceParser,
		)
		wantErr bool
	}{
		{
			name:    "Successfully save files",
			repoStr: "github.com/owner/repo",
			setupMock: func(t *testing.T, ctrl *gomock.Controller, files []string) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				tempDir := t.TempDir()
				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(tempDir, "github.com", "owner", "repo")

				// Create test file
				testFile := filepath.Join(tempDir, ".gitignore")
				if err := os.WriteFile(testFile, []byte("*.log\n"), 0o644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				files[0] = testFile

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// Overlay and hook creation
				overlayID := uuid.New()
				overlayService.EXPECT().Add(ctx, gomock.Any()).Return(overlayID.String(), nil)
				mockOverlay := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay.EXPECT().UUID().Return(overlayID)
				overlayService.EXPECT().Get(ctx, overlayID.String()).Return(mockOverlay, nil)

				hookID := uuid.New().String()
				hs.EXPECT().Add(ctx, gomock.Any()).Return(hookID, nil)

				// Save auto extra
				es.EXPECT().AddAutoExtra(ctx, ref, ref, gomock.Any()).Return(uuid.New().String(), nil)

				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			files:   []string{"placeholder"},
			wantErr: false,
		},
		{
			name:    "No files provided",
			repoStr: "github.com/owner/repo",
			files:   []string{},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, files []string) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				tempDir := t.TempDir()
				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(tempDir, "github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
		{
			name:    "Auto extra already exists",
			repoStr: "github.com/owner/repo",
			files:   []string{"test.txt"},
			setupMock: func(t *testing.T, ctrl *gomock.Controller, files []string) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*git_mock.MockGitService,
				*overlay_mock.MockOverlayService,
				*script_mock.MockScriptService,
				*hook_mock.MockHookService,
				*extra_mock.MockExtraService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				gs := git_mock.NewMockGitService(ctrl)
				overlayService := overlay_mock.NewMockOverlayService(ctrl)
				hs := hook_mock.NewMockHookService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				tempDir := t.TempDir()
				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(tempDir, "github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

				// Auto extra already exists
				existingExtra := extra.NewAutoExtra(
					uuid.New().String(),
					ref,
					ref,
					[]extra.Item{},
					time.Now(),
				)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(existingExtra, nil)

				return ws, fs, gs, overlayService, ss, hs, es, rp
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ws, fs, gs, overlayService, ss, hs, es, rp := tc.setupMock(t, ctrl, tc.files)
			uc := testtarget.NewUsecase(ws, fs, gs, overlayService, ss, hs, es, rp)

			err := uc.SaveFiles(ctx, tc.repoStr, tc.files)
			if (err != nil) != tc.wantErr {
				t.Errorf("SaveFiles() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
