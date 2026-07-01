package clone

import (
	"context"
	"errors"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	// Define test cases
	tests := []struct {
		name         string
		refWithAlias string
		setupMocks   func(
			mockHosting *hosting_mock.MockHostingService,
			mockWorkspace *workspace_mock.MockWorkspaceService,
			mockFinder *workspace_mock.MockFinderService,
			mockLayout *workspace_mock.MockLayoutService,
			mockOverlay *overlay_mock.MockOverlayService,
			mockScript *script_mock.MockScriptService,
			mockHook *hook_mock.MockHookService,
			mockRefParser *repository_mock.MockReferenceParser,
			mockGit *git_mock.MockGitService,
		)
		errorContains string
	}{
		{
			name:         "Success: Clone repository normally",
			refWithAlias: "github.com/kyoh86/gogh",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				tmpDir := t.TempDir()
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/gogh").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil).
					AnyTimes()

				pseudoCloneURL := "https://github.com/kyoh86/gogh.git"
				// Get repository information
				mockHosting.EXPECT().
					GetRepository(gomock.Any(), repository.NewReference("github.com", "kyoh86", "gogh")).
					Return(&hosting.Repository{
						Ref:         ref,
						URL:         "https://github.com/kyoh86/gogh",
						CloneURL:    pseudoCloneURL,
						UpdatedAt:   time.Now(),
						Description: "GoGH - GO Github Handler",
						Homepage:    "https://github.com/kyoh86/gogh",
						Language:    "Go",
					}, nil)
				mockHosting.EXPECT().
					GetTokenFor(gomock.Any(), "github.com", "kyoh86").Return("kyoh86", auth.Token{}, nil)

				pseudoPath := filepath.Join(tmpDir, "github.com/kyoh86/gogh")
				// Get primary layout to clone into
				mockLayout.EXPECT().PathFor(ref).Return(pseudoPath)
				mockWorkspace.EXPECT().GetPrimaryLayout().Return(mockLayout)

				// Verify that try.Execute is called
				// Since we're calling the actual function instead of a mock,
				// set expectations for GitService.Clone
				mockGit.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "kyoh86", "").Return(mockGit, nil)
				mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				// New bare + worktree flow
				mockGit.EXPECT().EnsureRemoteFetchRefspec(gomock.Any(), pseudoPath, "origin").Return(nil)
				mockGit.EXPECT().Fetch(gomock.Any(), pseudoPath, "origin").Return(nil)
				mockGit.EXPECT().SetRemoteHead(gomock.Any(), pseudoPath, "origin").Return(nil)
				mockGit.EXPECT().CreateBranch(gomock.Any(), pseudoPath, "main", "origin/HEAD").Return(nil)
				mockGit.EXPECT().AddWorktree(gomock.Any(), pseudoPath, "main", ".wt/main").Return(nil)
				mockGit.EXPECT().SetDefaultRemotes(gomock.Any(), pseudoPath, []string{pseudoCloneURL}).Return(nil)

				// Hook finds the repository by reference
				mockFinder.EXPECT().
					FindByReference(gomock.Any(), gomock.Any(), ref).
					Return(repository.NewLocation(
						pseudoPath,
						"github.com",
						"kyoh86",
						"gogh",
					), nil)
				overlayID := uuid.New()
				scriptID := uuid.New()
				mockHook.EXPECT().
					ListFor(ref, hook.EventPostClone).Return(func(yield func(hook.Hook, error) bool) {
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-clone-example",
						OperationType: hook.OperationTypeOverlay,
						OperationID:   overlayID,
					}), nil) {
						return
					}
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-clone-example",
						OperationType: hook.OperationTypeScript,
						OperationID:   scriptID,
					}), nil) {
						return
					}
				})
				mockOverlay.EXPECT().
					Get(gomock.Any(), overlayID.String()).Return(overlay.NewOverlay(overlay.Entry{
					Name:         "example-overlay",
					RelativePath: "path/to/overlay",
				}), nil)
				mockOverlay.EXPECT().
					Open(gomock.Any(), overlayID.String()).Return(io.NopCloser(strings.NewReader("overlay content")), nil)
				// Script cannot be run in this test, so we just return error
				mockScript.EXPECT().
					Open(gomock.Any(), scriptID.String()).Return(nil, errors.New("script error"))
			},
			errorContains: "script error",
		},
		{
			name:         "Error: Reference parsing error",
			refWithAlias: "invalid-reference",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				// Reference parsing error
				mockRefParser.EXPECT().
					ParseWithAlias("invalid-reference").
					Return(&repository.ReferenceWithAlias{}, errors.New("invalid reference format"))
			},
			errorContains: "invalid reference format",
		},
		{
			name:         "Success: Direct clone fallback when repository metadata fetch fails",
			refWithAlias: "github.com/kyoh86/not-found",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				tmpDir := t.TempDir()
				ref := repository.NewReference("github.com", "kyoh86", "not-found")
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/not-found").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil).
					Times(2)

				mockHosting.EXPECT().
					GetRepository(gomock.Any(), ref).
					Return(&hosting.Repository{}, errors.New("repository not found"))
				mockHosting.EXPECT().
					GetURLOf(ref).
					Return(mustParseURL(t, "https://github.com/kyoh86/not-found"), nil)

				pseudoCloneURL := "https://github.com/kyoh86/not-found.git"
				pseudoPath := filepath.Join(tmpDir, "github.com/kyoh86/not-found")
				mockLayout.EXPECT().PathFor(ref).Return(pseudoPath)
				mockWorkspace.EXPECT().GetPrimaryLayout().Return(mockLayout)

				mockGit.EXPECT().
					Clone(gomock.Any(), pseudoCloneURL, pseudoPath, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, _ string, opts git.CloneOptions) error {
						if !opts.IsBare {
							t.Fatalf("expected bare clone for worktree mode")
						}
						if !opts.UseSystemGit {
							t.Fatalf("expected direct fallback to use system git")
						}
						return nil
					})
				mockGit.EXPECT().AddWorktree(gomock.Any(), pseudoPath, "main", ".wt/main").Return(nil)
				mockGit.EXPECT().SetDefaultRemotes(gomock.Any(), pseudoPath, []string{pseudoCloneURL}).Return(nil)

				mockFinder.EXPECT().
					FindByReference(gomock.Any(), gomock.Any(), ref).
					Return(repository.NewLocation(
						pseudoPath,
						"github.com",
						"kyoh86",
						"not-found",
					), nil)
				mockHook.EXPECT().
					ListFor(ref, hook.EventPostClone).Return(func(yield func(hook.Hook, error) bool) {})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
			mockHosting := hosting_mock.NewMockHostingService(ctrl)
			mockWorkspace := workspace_mock.NewMockWorkspaceService(ctrl)
			mockFinder := workspace_mock.NewMockFinderService(ctrl)
			mockLayout := workspace_mock.NewMockLayoutService(ctrl)
			mockOverlay := overlay_mock.NewMockOverlayService(ctrl)
			mockScript := script_mock.NewMockScriptService(ctrl)
			mockHook := hook_mock.NewMockHookService(ctrl)
			mockRefParser := repository_mock.NewMockReferenceParser(ctrl)
			mockGit := git_mock.NewMockGitService(ctrl)

			// Setup mocks
			tt.setupMocks(mockHosting, mockWorkspace, mockFinder, mockLayout, mockOverlay, mockScript, mockHook, mockRefParser, mockGit)

			// Create Usecase to test
			usecase := NewUsecase(mockHosting, mockWorkspace, mockFinder, mockOverlay, mockScript, mockHook, mockRefParser, mockGit)

			// Execute test
			err := usecase.Execute(context.Background(), tt.refWithAlias, Options{
				TryCloneOptions: try.Options{
					Worktree: true,
				},
			})

			// Verify results
			switch {
			case tt.errorContains == "":
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			case err == nil:
				t.Errorf("Expected an error but got none")
			case !strings.Contains(err.Error(), tt.errorContains) && errors.Unwrap(err) == nil:
				t.Errorf("Expected error containing: %v, got: %v", tt.errorContains, err)
			}
		})
	}
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}
	return u
}
