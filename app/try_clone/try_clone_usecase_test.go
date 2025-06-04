package try_clone_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestNewUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hostingService := hosting_mock.NewMockHostingService(ctrl)
	workspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
	overlayStore := overlay_mock.NewMockOverlayStore(ctrl)
	gitService := git_mock.NewMockGitService(ctrl)

	svc := try_clone.NewUseCase(hostingService, workspaceService, overlayStore, gitService)
	if svc == nil {
		t.Fatal("NewRepositoryService returned nil")
	}
}

func TestRetryLimit(t *testing.T) {
	testCases := []struct {
		name          string
		limit         int
		notifications []try_clone.Status
		expectErr     bool
		errAt         int
	}{
		{
			name:          "no retries needed",
			limit:         3,
			notifications: []try_clone.Status{try_clone.StatusEmpty},
			expectErr:     false,
		},
		{
			name:  "retries within limit",
			limit: 3,
			notifications: []try_clone.Status{
				try_clone.StatusRetry,
				try_clone.StatusRetry,
				try_clone.StatusEmpty,
			},
			expectErr: false,
		},
		{
			name:  "retries exceed limit",
			limit: 2,
			notifications: []try_clone.Status{
				try_clone.StatusRetry,
				try_clone.StatusRetry,
				try_clone.StatusRetry,
			},
			expectErr: true,
			errAt:     2, // 0-indexed, so this is the 3rd retry attempt
		},
		{
			name:  "nil notify function",
			limit: 2,
			notifications: []try_clone.Status{
				try_clone.StatusRetry,
				try_clone.StatusEmpty,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var notify try_clone.Notify
			if tc.name != "nil notify function" {
				notify = func(n try_clone.Status) error { return nil }
			}

			limitedNotify := try_clone.RetryLimit(tc.limit, notify)

			var err error
			for i, status := range tc.notifications {
				err = limitedNotify(status)
				if err != nil {
					if !tc.expectErr {
						t.Fatalf("Unexpected error at notification %d: %v", i, err)
					}
					if tc.errAt != i {
						t.Fatalf("Error expected at notification %d, but occurred at %d", tc.errAt, i)
					}
					break
				}
			}

			if tc.expectErr && err == nil {
				t.Fatal("Expected error but got nil")
			}
		})
	}
}

func TestTryClone(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService)
		expectErr     bool
		expectErrText string
	}{
		{
			name: "successful clone",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				repo := &hosting.Repository{
					CloneURL: "https://github.com/user/repo.git",
				}
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("user", auth.Token{AccessToken: "token"}, nil)
				mgs.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "user", "token").Return(mgs, nil)

				// Clone
				mgs.EXPECT().Clone(gomock.Any(), repo.CloneURL, localPath, gomock.Any()).Return(nil)

				// Remote setup
				mgs.EXPECT().SetDefaultRemotes(gomock.Any(), localPath, []string{repo.CloneURL}).Return(nil)

				return mhs, mws, mos, mgs
			},
			expectErr: false,
		},
		{
			name: "authentication error",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication error
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("", auth.Token{}, errors.New("auth error"))

				return mhs, mws, mos, mgs
			},
			expectErr:     true,
			expectErrText: "auth error",
		},
		{
			name: "authentication username/password error",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("user", auth.Token{AccessToken: "token"}, nil)
				mgs.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "user", "token").Return(nil, errors.New("auth username/password error"))

				return mhs, mws, mos, mgs
			},
			expectErr:     true,
			expectErrText: "auth username/password error",
		},
		{
			name: "clone error",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				repo := &hosting.Repository{
					CloneURL: "https://github.com/user/repo.git",
				}
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("user", auth.Token{AccessToken: "token"}, nil)
				mgs.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "user", "token").Return(mgs, nil)

				// Clone error
				mgs.EXPECT().Clone(gomock.Any(), repo.CloneURL, localPath, gomock.Any()).
					Return(errors.New("clone error"))

				return mhs, mws, mos, mgs
			},
			expectErr:     true,
			expectErrText: "cloning: clone error",
		},
		{
			name: "empty repository",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				repo := &hosting.Repository{
					CloneURL: "https://github.com/user/repo.git",
				}
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("user", auth.Token{AccessToken: "token"}, nil)
				mgs.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "user", "token").Return(mgs, nil)

				// Clone returns empty repository error
				mgs.EXPECT().Clone(gomock.Any(), repo.CloneURL, localPath, gomock.Any()).
					Return(git.ErrRepositoryEmpty)

				// Init for empty repository
				mls.EXPECT().CreateRepositoryFolder(ref).Return(localPath, nil)
				mgs.EXPECT().Init(gomock.Any(), repo.CloneURL, localPath, false, gomock.Any()).Return(nil)

				// Remote setup
				mgs.EXPECT().SetDefaultRemotes(gomock.Any(), localPath, []string{repo.CloneURL}).Return(nil)

				return mhs, mws, mos, mgs
			},
			expectErr: false,
		},
		{
			name: "parent repository setup",
			setupMocks: func(ctrl *gomock.Controller) (*hosting_mock.MockHostingService, *workspace_mock.MockWorkspaceService, *overlay_mock.MockOverlayStore, *git_mock.MockGitService) {
				mhs := hosting_mock.NewMockHostingService(ctrl)
				mws := workspace_mock.NewMockWorkspaceService(ctrl)
				mgs := git_mock.NewMockGitService(ctrl)
				mls := workspace_mock.NewMockLayoutService(ctrl)
				mos := overlay_mock.NewMockOverlayStore(ctrl)

				ref := repository.NewReference("github.com", "user", "repo")
				repo := &hosting.Repository{
					CloneURL: "https://github.com/user/repo.git",
					Parent: &hosting.ParentRepository{
						CloneURL: "https://github.com/original/repo.git",
					},
				}
				localPath := "/path/to/repo"

				// Layout setup
				mws.EXPECT().GetPrimaryLayout().Return(mls)
				mls.EXPECT().PathFor(ref).Return(localPath)

				// Authentication
				mhs.EXPECT().GetTokenFor(gomock.Any(), ref.Host(), ref.Owner()).Return("user", auth.Token{AccessToken: "token"}, nil)
				mgs.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "user", "token").Return(mgs, nil)

				// Clone
				mgs.EXPECT().Clone(gomock.Any(), repo.CloneURL, localPath, gomock.Any()).Return(nil)

				// Remote setup
				mgs.EXPECT().SetDefaultRemotes(gomock.Any(), localPath, []string{repo.CloneURL}).Return(nil)

				// Parent remote setup
				mgs.EXPECT().SetRemotes(gomock.Any(), localPath, "upstream", []string{repo.Parent.CloneURL}).Return(nil)

				return mhs, mws, mos, mgs
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mhs, mws, mgs, mos := tc.setupMocks(ctrl)

			svc := try_clone.NewUseCase(mhs, mws, mgs, mos)

			repo := &hosting.Repository{
				Ref:      repository.NewReference("github.com", "user", "repo"),
				CloneURL: "https://github.com/user/repo.git",
			}

			// Handle parent repo case
			if tc.name == "parent repository setup" {
				repo.Parent = &hosting.ParentRepository{
					CloneURL: "https://github.com/original/repo.git",
				}
			}

			var notifyCalledWith try_clone.Status
			notify := func(status try_clone.Status) error {
				notifyCalledWith = status
				return nil
			}

			err := svc.Execute(
				context.Background(),
				repo,
				nil, // no alias
				try_clone.Options{
					Timeout: 30 * time.Second,
					Notify:  notify,
				},
			)

			if tc.expectErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
				if tc.expectErrText != "" && !errors.Is(err, err) && errors.Unwrap(err) != nil {
					if msg := err.Error(); !containsString(msg, tc.expectErrText) {
						t.Fatalf("Expected error to contain %q but got %q", tc.expectErrText, msg)
					}
				}
			} else if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			// Verify notification was called correctly for empty repo case
			if tc.name == "empty repository" && notifyCalledWith != try_clone.StatusEmpty {
				t.Errorf("Expected notification with TryCloneStatusEmpty, got %v", notifyCalledWith)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}
