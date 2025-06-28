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

				// Parse reference
				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)

				// Find repository
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

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

				ref := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					"/path/to/repo",
					"github.com",
					"owner",
					"repo",
				)

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)

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

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

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

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// List excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, ".gitignore"), nil)
						yield(filepath.Join(tempDir, "config.toml"), nil)
					},
				)

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

				rp.EXPECT().Parse("github.com/owner/repo").Return(&ref, nil)
				fs.EXPECT().FindByReference(ctx, ws, ref).Return(location, nil)
				es.EXPECT().GetAutoExtra(ctx, ref).Return(nil, errors.New("not found"))

				// List excluded files
				gs.EXPECT().ListExcludedFiles(ctx, tempDir, nil).Return(
					func(yield func(string, error) bool) {
						yield(filepath.Join(tempDir, ".gitignore"), nil)
					},
				)

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
