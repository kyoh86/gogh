package create_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/create"
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
	// Define test cases
	tests := []struct {
		name         string
		refWithAlias string
		options      testtarget.Options
		setupMocks   func(
			mockHosting *hosting_mock.MockHostingService,
			mockWorkspace *workspace_mock.MockWorkspaceService,
			mockFinder *workspace_mock.MockFinderService,
			mockLayout *workspace_mock.MockLayoutService,
			mockOverlay *overlay_mock.MockOverlayService,
			mockHook *hook_mock.MockHookService,
			mockRefParser *repository_mock.MockReferenceParser,
			mockGit *git_mock.MockGitService,
		)
		expectedError bool
		errorContains string
	}{
		{
			name:         "Success: Create and clone repository",
			refWithAlias: "github.com/kyoh86/new-repo",
			options: testtarget.Options{
				TryCloneOptions: try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryOptions{
					Description: "New repository",
					Private:     true,
				},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				ref := repository.NewReference("github.com", "kyoh86", "new-repo")
				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/new-repo").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil).AnyTimes()

				pseudoCloneURL := "https://github.com/kyoh86/new-repo.git"
				// Create repository
				mockHosting.EXPECT().
					CreateRepository(gomock.Any(), ref, hosting.CreateRepositoryOptions{
						Description: "New repository",
						Private:     true,
					}).
					Return(&hosting.Repository{
						Ref:         ref,
						URL:         "https://github.com/kyoh86/new-repo",
						CloneURL:    pseudoCloneURL,
						UpdatedAt:   time.Now(),
						Description: "New repository",
						Homepage:    "https://github.com/kyoh86/new-repo",
						Language:    "Go",
					}, nil)
				mockHosting.EXPECT().
					GetTokenFor(gomock.Any(), "github.com", "kyoh86").Return("kyoh86", auth.Token{}, nil)

				pseudoPath := "/path/to/workspace/github.com/kyoh86/new-repo"
				// Get primary layout to clone into
				mockLayout.EXPECT().PathFor(ref).Return(pseudoPath)
				mockWorkspace.EXPECT().GetPrimaryLayout().Return(mockLayout)

				// Set expectations for GitService.Clone
				mockGit.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "kyoh86", "").Return(mockGit, nil)
				mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockGit.EXPECT().SetDefaultRemotes(gomock.Any(), pseudoPath, []string{pseudoCloneURL}).Return(nil)

				// Overlay application
				mockOverlay.EXPECT().
					ListOverlays().Return(func(yield func(*overlay.Overlay, error) bool) {})
				// Hook application
				mockHook.EXPECT().
					ListHooks().Return(func(yield func(*hook.Hook, error) bool) {}).Times(2)
			},
			expectedError: false,
		},
		{
			name:         "Error: Reference parsing error",
			refWithAlias: "invalid-reference",
			options: testtarget.Options{
				TryCloneOptions:   try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryOptions{},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				// Reference parsing error
				mockRefParser.EXPECT().
					ParseWithAlias("invalid-reference").
					Return(&repository.ReferenceWithAlias{}, errors.New("invalid reference format"))
			},
			expectedError: true,
			errorContains: "invalid ref: invalid reference format",
		},
		{
			name:         "Error: Repository creation error",
			refWithAlias: "github.com/kyoh86/repo-creation-failed",
			options: testtarget.Options{
				TryCloneOptions:   try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryOptions{},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				ref := repository.NewReference("github.com", "kyoh86", "repo-creation-failed")
				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/repo-creation-failed").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil)

				// Repository creation error
				mockHosting.EXPECT().
					CreateRepository(gomock.Any(), ref, hosting.CreateRepositoryOptions{}).
					Return(&hosting.Repository{}, errors.New("failed to create repository"))
			},
			expectedError: true,
			errorContains: "creating: failed to create repository",
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
			mockHook := hook_mock.NewMockHookService(ctrl)
			mockRefParser := repository_mock.NewMockReferenceParser(ctrl)
			mockGit := git_mock.NewMockGitService(ctrl)

			// Setup mocks
			tt.setupMocks(mockHosting, mockWorkspace, mockFinder, mockLayout, mockOverlay, mockHook, mockRefParser, mockGit)

			// Create UseCase to test
			useCase := testtarget.NewUseCase(mockHosting, mockWorkspace, mockFinder, mockOverlay, mockHook, mockRefParser, mockGit)

			// Execute test
			err := useCase.Execute(context.Background(), tt.refWithAlias, tt.options)

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing: %v, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
