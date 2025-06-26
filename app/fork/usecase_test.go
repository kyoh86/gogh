package fork_test

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/core/auth"
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
	testCases := []struct {
		name       string
		source     string
		target     string
		setupMocks func(
			mockHosting *hosting_mock.MockHostingService,
			mockWorkspace *workspace_mock.MockWorkspaceService,
			mockFinder *workspace_mock.MockFinderService,
			mockLayout *workspace_mock.MockLayoutService,
			mockOverlay *overlay_mock.MockOverlayService,
			mockScript *script_mock.MockScriptService,
			mockHook *hook_mock.MockHookService,
			mockDefaultName *repository_mock.MockDefaultNameService,
			mockReferenceParser *repository_mock.MockReferenceParser,
			mockGit *git_mock.MockGitService,
		)
		expectErrText string
	}{
		{
			name:   "successful fork and clone",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				tmpDir := t.TempDir()
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mockReferenceParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockReferenceParser.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil).AnyTimes()

				// Fork repository
				forkedRepo := &hosting.Repository{Ref: targetRef}
				mockHosting.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, targetRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository
				mockHosting.EXPECT().
					GetTokenFor(gomock.Any(), targetRef.Host(), targetRef.Owner()).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mockGit.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mockGit, nil)
				mockGit.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockGit.EXPECT().
					SetDefaultRemotes(gomock.Any(), gomock.Any(), []string{forkedRepo.CloneURL}).
					Return(nil)

				pseudoPath := filepath.Join(tmpDir, "github.com/target/repo")
				mockLayout.EXPECT().
					PathFor(targetRef).
					Return(pseudoPath)
				mockWorkspace.EXPECT().
					GetPrimaryLayout().
					Return(mockLayout)

				// Hook finds the repository by reference
				mockFinder.EXPECT().
					FindByReference(gomock.Any(), gomock.Any(), targetRef).
					Return(repository.NewLocation(
						pseudoPath,
						"github.com",
						"target",
						"repo",
					), nil)
				mockHook.EXPECT().
					ListFor(targetRef, hook.EventPostFork).Return(func(yield func(hook.Hook, error) bool) {
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-fork-example",
						OperationType: hook.OperationTypeOverlay,
						OperationID:   "overlay-id",
					}), nil) {
						return
					}
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-fork-example",
						OperationType: hook.OperationTypeScript,
						OperationID:   "script-id",
					}), nil) {
						return
					}
				})
				mockOverlay.EXPECT().
					Get(gomock.Any(), "overlay-id").Return(overlay.NewOverlay(overlay.Entry{
					Name:         "example-overlay",
					RelativePath: "path/to/overlay",
				}), nil)
				mockOverlay.EXPECT().
					Open(gomock.Any(), "overlay-id").Return(io.NopCloser(strings.NewReader("overlay content")), nil)
				// Script cannot be run in this test, so we just return error
				mockScript.EXPECT().
					Open(gomock.Any(), "script-id").Return(nil, errors.New("script error"))
			},
			expectErrText: "script error",
		},
		{
			name:   "invalid source reference",
			source: "invalid-source",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				mockReferenceParser.EXPECT().
					Parse("invalid-source").
					Return(nil, errors.New("invalid source reference"))
			},
			expectErrText: "invalid source",
		},
		{
			name:   "empty target with default owner",
			source: "github.com/source/repo",
			target: "",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				tmpDir := t.TempDir()
				sourceRef := repository.NewReference("github.com", "source", "repo")
				mockReferenceParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)

				mockDefaultName.EXPECT().
					GetDefaultOwnerFor("github.com").
					Return("default-owner", nil)

				defaultRef := repository.NewReference("github.com", "default-owner", "repo")
				mockReferenceParser.EXPECT().
					ParseWithAlias("github.com/default-owner/repo").
					Return(&repository.ReferenceWithAlias{Reference: defaultRef}, nil).AnyTimes()

				// Fork repository
				forkedRepo := &hosting.Repository{
					Ref: repository.NewReference("github.com", "default-owner", "repo"),
				}
				mockHosting.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, defaultRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository
				mockHosting.EXPECT().
					GetTokenFor(gomock.Any(), defaultRef.Host(), defaultRef.Owner()).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mockGit.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mockGit, nil)
				mockGit.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mockGit.EXPECT().
					SetDefaultRemotes(gomock.Any(), gomock.Any(), []string{forkedRepo.CloneURL}).
					Return(nil)

				pseudoPath := filepath.Join(tmpDir, "github.com/target/repo")
				mockLayout.EXPECT().
					PathFor(defaultRef).
					Return(pseudoPath)
				mockWorkspace.EXPECT().
					GetPrimaryLayout().
					Return(mockLayout)

				// Hook finds the repository by reference
				mockFinder.EXPECT().
					FindByReference(gomock.Any(), gomock.Any(), defaultRef).
					Return(repository.NewLocation(
						pseudoPath,
						"github.com",
						"default-owner",
						"repo",
					), nil)
				mockHook.EXPECT().
					ListFor(defaultRef, hook.EventPostFork).Return(func(yield func(hook.Hook, error) bool) {
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-fork-example",
						OperationType: hook.OperationTypeOverlay,
						OperationID:   "overlay-id",
					}), nil) {
						return
					}
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-fork-example",
						OperationType: hook.OperationTypeScript,
						OperationID:   "script-id",
					}), nil) {
						return
					}
				})
				mockOverlay.EXPECT().
					Get(gomock.Any(), "overlay-id").Return(overlay.NewOverlay(overlay.Entry{
					Name:         "example-overlay",
					RelativePath: "path/to/overlay",
				}), nil)
				mockOverlay.EXPECT().
					Open(gomock.Any(), "overlay-id").Return(io.NopCloser(strings.NewReader("overlay content")), nil)
				// Script cannot be run in this test, so we just return error
				mockScript.EXPECT().
					Open(gomock.Any(), "script-id").Return(nil, errors.New("script error"))
			},
			expectErrText: "script error",
		},
		{
			name:   "fork error",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mockReferenceParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockReferenceParser.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil)

				// Fork repository
				mockHosting.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, targetRef, gomock.Any()).
					Return(nil, errors.New("fork error"))
			},
			expectErrText: "requesting fork",
		},
		{
			name:   "clone error",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mockReferenceParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockReferenceParser.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil)

				// Fork repository
				forkedRepo := &hosting.Repository{
					Ref: repository.NewReference("github.com", "target", "repo"),
				}
				mockHosting.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, targetRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository error
				mockHosting.EXPECT().
					GetTokenFor(gomock.Any(), targetRef.Host(), targetRef.Owner()).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mockGit.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mockGit, nil)
				mockGit.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("cloning error"))

				mockLayout.EXPECT().
					PathFor(targetRef).
					Return("/path/to/repo")
				mockWorkspace.EXPECT().
					GetPrimaryLayout().
					Return(mockLayout)
			},
			expectErrText: "cloning forked repository",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockHostingService := hosting_mock.NewMockHostingService(ctrl)
				mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
				mockFinderService := workspace_mock.NewMockFinderService(ctrl)
				mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
				mockScriptService := script_mock.NewMockScriptService(ctrl)
				mockHookService := hook_mock.NewMockHookService(ctrl)
				mockLayoutService := workspace_mock.NewMockLayoutService(ctrl)
				mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
				mockReferenceParser := repository_mock.NewMockReferenceParser(ctrl)
				mockGitService := git_mock.NewMockGitService(ctrl)

				tc.setupMocks(
					mockHostingService,
					mockWorkspaceService,
					mockFinderService,
					mockLayoutService,
					mockOverlayService,
					mockScriptService,
					mockHookService,
					mockDefaultNameService,
					mockReferenceParser,
					mockGitService,
				)

				usecase := fork.NewUsecase(
					mockHostingService,
					mockWorkspaceService,
					mockFinderService,
					mockOverlayService,
					mockScriptService,
					mockHookService,
					mockDefaultNameService,
					mockReferenceParser,
					mockGitService,
				)

				opts := fork.Options{
					TryCloneOptions: try.Options{
						Timeout: 30 * time.Second,
						Notify:  func(msg try.Status) error { return nil },
					},
					Target: tc.target,
				}

				err := usecase.Execute(context.Background(), tc.source, opts)

				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.expectErrText != "" && !errors.Is(err, err) && errors.Unwrap(err) != nil {
					if msg := err.Error(); !containsString(msg, tc.expectErrText) {
						t.Fatalf("Expected error to contain %q but got %q", tc.expectErrText, msg)
					}
				}
			},
		)
	}
}

func containsString(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}
