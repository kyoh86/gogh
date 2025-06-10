package fork_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/app/try_clone"
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
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
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
			mockHook *hook_mock.MockHookService,
			mockDefaultName *repository_mock.MockDefaultNameService,
			mockReferenceParser *repository_mock.MockReferenceParser,
			mockGit *git_mock.MockGitService,
		)
		expectErr     bool
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

				mockLayout.EXPECT().
					PathFor(targetRef).
					Return("/path/to/repo")
				mockWorkspace.EXPECT().
					GetPrimaryLayout().
					Return(mockLayout)

				// Overlay application
				mockOverlay.EXPECT().
					ListOverlays().Return(func(yield func(*overlay.Overlay, error) bool) {})
				// Hook application
				mockHook.EXPECT().
					ListHooks().Return(func(yield func(*hook.Hook, error) bool) {}).Times(2)
			},
			expectErr: false,
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
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				mockReferenceParser.EXPECT().
					Parse("invalid-source").
					Return(nil, errors.New("invalid source reference"))
			},
			expectErr:     true,
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
				mockHook *hook_mock.MockHookService,
				mockDefaultName *repository_mock.MockDefaultNameService,
				mockReferenceParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
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

				mockLayout.EXPECT().
					PathFor(defaultRef).
					Return("/path/to/repo")
				mockWorkspace.EXPECT().
					GetPrimaryLayout().
					Return(mockLayout)

				// Overlay application
				mockOverlay.EXPECT().
					ListOverlays().Return(func(yield func(*overlay.Overlay, error) bool) {})
				// Hook application
				mockHook.EXPECT().
					ListHooks().Return(func(yield func(*hook.Hook, error) bool) {}).Times(2)
			},
			expectErr: false,
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
			expectErr:     true,
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
			expectErr:     true,
			expectErrText: "cloning forked repository",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHostingService := hosting_mock.NewMockHostingService(ctrl)
			mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
			mockFinderService := workspace_mock.NewMockFinderService(ctrl)
			mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
			mockHookService := hook_mock.NewMockHookService(ctrl)
			mockLayoutService := workspace_mock.NewMockLayoutService(ctrl)
			mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
			mockReferenceParser := repository_mock.NewMockReferenceParser(ctrl)
			mockGitService := git_mock.NewMockGitService(ctrl)

			tc.setupMocks(mockHostingService, mockWorkspaceService, mockFinderService, mockLayoutService, mockOverlayService, mockHookService, mockDefaultNameService, mockReferenceParser, mockGitService)

			useCase := fork.NewUseCase(
				mockHostingService,
				mockWorkspaceService,
				mockFinderService,
				mockOverlayService,
				mockHookService,
				mockDefaultNameService,
				mockReferenceParser,
				mockGitService,
			)

			opts := fork.Options{
				TryCloneOptions: try_clone.Options{
					Timeout: 30 * time.Second,
					Notify:  func(msg try_clone.Status) error { return nil },
				},
				Target: tc.target,
			}

			err := useCase.Execute(context.Background(), tc.source, opts)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.expectErrText != "" && !errors.Is(err, err) && errors.Unwrap(err) != nil {
					if msg := err.Error(); !containsString(msg, tc.expectErrText) {
						t.Fatalf("Expected error to contain %q but got %q", tc.expectErrText, msg)
					}
				}
			} else if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}
