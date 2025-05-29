package fork_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	testCases := []struct {
		name          string
		source        string
		target        string
		setupMocks    func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService)
		expectErr     bool
		expectErrText string
	}{
		{
			name:   "successful fork and clone",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mrp.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mrp.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil)

				// Fork repository
				forkedRepo := &hosting.Repository{Ref: targetRef}
				mhs.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, targetRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository
				mhs.EXPECT().
					GetTokenFor(gomock.Any(), targetRef).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mgs.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mgs, nil)
				mgs.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mgs.EXPECT().
					SetDefaultRemotes(gomock.Any(), gomock.Any(), []string{forkedRepo.CloneURL}).
					Return(nil)

				mls.EXPECT().
					PathFor(targetRef).
					Return("/path/to/repo")
				mws.EXPECT().
					GetPrimaryLayout().
					Return(mls)
			},
			expectErr: false,
		},
		{
			name:   "invalid source reference",
			source: "invalid-source",
			setupMocks: func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService) {
				mrp.EXPECT().
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
			setupMocks: func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				mrp.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)

				mdns.EXPECT().
					GetDefaultOwnerFor("github.com").
					Return("default-owner", nil)

				defaultRef := repository.NewReference("github.com", "default-owner", "repo")

				// Fork repository
				forkedRepo := &hosting.Repository{
					Ref: repository.NewReference("github.com", "default-owner", "repo"),
				}
				mhs.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, defaultRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository
				mhs.EXPECT().
					GetTokenFor(gomock.Any(), defaultRef).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mgs.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mgs, nil)
				mgs.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mgs.EXPECT().
					SetDefaultRemotes(gomock.Any(), gomock.Any(), []string{forkedRepo.CloneURL}).
					Return(nil)

				mls.EXPECT().
					PathFor(defaultRef).
					Return("/path/to/repo")
				mws.EXPECT().
					GetPrimaryLayout().
					Return(mls)
			},
			expectErr: false,
		},
		{
			name:   "fork error",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mrp.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mrp.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil)

				// Fork repository
				mhs.EXPECT().
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
			setupMocks: func(mhs *hosting_mock.MockHostingService, mws *workspace_mock.MockWorkspaceService, mls *workspace_mock.MockLayoutService, mdns *repository_mock.MockDefaultNameService, mrp *repository_mock.MockReferenceParser, mgs *git_mock.MockGitService) {
				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "target", "repo")
				targetRefWithAlias := &repository.ReferenceWithAlias{
					Reference: targetRef,
				}

				// Parse references
				mrp.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mrp.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRefWithAlias, nil)

				// Fork repository
				forkedRepo := &hosting.Repository{
					Ref: repository.NewReference("github.com", "target", "repo"),
				}
				mhs.EXPECT().
					ForkRepository(gomock.Any(), sourceRef, targetRef, gomock.Any()).
					Return(forkedRepo, nil)

				// Clone repository error
				mhs.EXPECT().
					GetTokenFor(gomock.Any(), targetRef).
					Return("target-auth-user", auth.Token{AccessToken: "target-auth-token"}, nil) // Get token for target repository to clone
				mgs.EXPECT().
					AuthenticateWithUsernamePassword(gomock.Any(), "target-auth-user", "target-auth-token").
					Return(mgs, nil)
				mgs.EXPECT().
					Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("cloning error"))

				mls.EXPECT().
					PathFor(targetRef).
					Return("/path/to/repo")
				mws.EXPECT().
					GetPrimaryLayout().
					Return(mls)
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
			mockLayoutService := workspace_mock.NewMockLayoutService(ctrl)
			mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
			mockReferenceParser := repository_mock.NewMockReferenceParser(ctrl)
			mockGitService := git_mock.NewMockGitService(ctrl)

			tc.setupMocks(mockHostingService, mockWorkspaceService, mockLayoutService, mockDefaultNameService, mockReferenceParser, mockGitService)

			useCase := fork.NewUseCase(
				mockHostingService,
				mockWorkspaceService,
				mockDefaultNameService,
				mockReferenceParser,
				mockGitService,
			)

			opts := fork.Options{
				RequestTimeout: 30 * time.Second,
				TryCloneNotify: func(msg service.TryCloneStatus) error { return nil },
				Target:         tc.target,
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
