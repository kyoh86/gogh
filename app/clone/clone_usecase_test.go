package clone

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

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
			name:         "Success: Clone repository normally",
			refWithAlias: "github.com/kyoh86/gogh",
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

				pseudoPath := "/path/to/workspace/github.com/kyoh86/gogh"
				// Get primary layout to clone into
				mockLayout.EXPECT().PathFor(ref).Return(pseudoPath)
				mockWorkspace.EXPECT().GetPrimaryLayout().Return(mockLayout)

				// Verify that try_clone.Execute is called
				// Since we're calling the actual function instead of a mock,
				// set expectations for GitService.Clone
				mockGit.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "kyoh86", "").Return(mockGit, nil)
				mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockGit.EXPECT().SetDefaultRemotes(gomock.Any(), pseudoPath, []string{pseudoCloneURL}).Return(nil)

				// Overlay application
				mockOverlay.EXPECT().
					List().Return(func(yield func(*overlay.Overlay, error) bool) {})
				mockHook.EXPECT().
					List().Return(func(yield func(*hook.Hook, error) bool) {}).Times(2)
			},
			expectedError: false,
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
			errorContains: "invalid reference format",
		},
		{
			name:         "Error: Repository fetch error",
			refWithAlias: "github.com/kyoh86/not-found",
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
				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/not-found").
					Return(&repository.ReferenceWithAlias{
						Reference: repository.NewReference("github.com", "kyoh86", "not-found"),
						Alias:     nil,
					}, nil)

				// Repository fetch error
				mockHosting.EXPECT().
					GetRepository(gomock.Any(), repository.NewReference("github.com", "kyoh86", "not-found")).
					Return(&hosting.Repository{}, errors.New("repository not found"))
			},
			expectedError: true,
			errorContains: "repository not found",
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
			useCase := NewUseCase(mockHosting, mockWorkspace, mockFinder, mockOverlay, mockHook, mockRefParser, mockGit)

			// Execute test
			err := useCase.Execute(context.Background(), tt.refWithAlias, Options{
				TryCloneOptions: try_clone.Options{},
			})

			// Verify results
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) && errors.Unwrap(err) == nil {
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
