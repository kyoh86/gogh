package create_from_template_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/create_from_template"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
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
		template     repository.Reference
		options      testtarget.CreateFromTemplateOptions
		setupMocks   func(
			mockHosting *hosting_mock.MockHostingService,
			mockWorkspace *workspace_mock.MockWorkspaceService,
			mockLayout *workspace_mock.MockLayoutService,
			mockOverlay *overlay_mock.MockOverlayService,
			mockRefParser *repository_mock.MockReferenceParser,
			mockGit *git_mock.MockGitService,
		)
		expectedError bool
		errorContains string
	}{
		{
			name:         "Success: Create from template and clone repository",
			refWithAlias: "github.com/kyoh86/new-repo",
			template:     repository.NewReference("github.com", "kyoh86", "template-repo"),
			options: testtarget.CreateFromTemplateOptions{
				TryCloneOptions: try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryFromTemplateOptions{
					Description: "New repository from template",
					Private:     true,
				},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				ref := repository.NewReference("github.com", "kyoh86", "new-repo")
				templateRef := repository.NewReference("github.com", "kyoh86", "template-repo")

				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/new-repo").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil)

				pseudoCloneURL := "https://github.com/kyoh86/new-repo.git"
				// Create repository from template
				mockHosting.EXPECT().
					CreateRepositoryFromTemplate(gomock.Any(), ref, templateRef, hosting.CreateRepositoryFromTemplateOptions{
						Description: "New repository from template",
						Private:     true,
					}).
					Return(&hosting.Repository{
						Ref:         ref,
						URL:         "https://github.com/kyoh86/new-repo",
						CloneURL:    pseudoCloneURL,
						UpdatedAt:   time.Now(),
						Description: "New repository from template",
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
			},
			expectedError: false,
		},
		{
			name:         "Error: Reference parsing error",
			refWithAlias: "invalid-reference",
			template:     repository.NewReference("github.com", "kyoh86", "template-repo"),
			options: testtarget.CreateFromTemplateOptions{
				TryCloneOptions:   try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryFromTemplateOptions{},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				// Reference parsing error
				mockRefParser.EXPECT().
					ParseWithAlias("invalid-reference").
					Return(&repository.ReferenceWithAlias{}, errors.New("invalid reference format"))
			},
			expectedError: true,
			errorContains: "invalid reference: invalid reference format",
		},
		{
			name:         "Error: Repository creation from template error",
			refWithAlias: "github.com/kyoh86/repo-creation-failed",
			template:     repository.NewReference("github.com", "kyoh86", "template-repo"),
			options: testtarget.CreateFromTemplateOptions{
				TryCloneOptions:   try_clone.Options{},
				RepositoryOptions: hosting.CreateRepositoryFromTemplateOptions{},
			},
			setupMocks: func(
				mockHosting *hosting_mock.MockHostingService,
				mockWorkspace *workspace_mock.MockWorkspaceService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				ref := repository.NewReference("github.com", "kyoh86", "repo-creation-failed")
				templateRef := repository.NewReference("github.com", "kyoh86", "template-repo")

				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/repo-creation-failed").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil)

				// Repository creation from template error
				mockHosting.EXPECT().
					CreateRepositoryFromTemplate(gomock.Any(), ref, templateRef, hosting.CreateRepositoryFromTemplateOptions{}).
					Return(&hosting.Repository{}, errors.New("failed to create repository from template"))
			},
			expectedError: true,
			errorContains: "creating repository from template: failed to create repository from template",
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
			mockLayout := workspace_mock.NewMockLayoutService(ctrl)
			mockOverlay := overlay_mock.NewMockOverlayService(ctrl)
			mockRefParser := repository_mock.NewMockReferenceParser(ctrl)
			mockGit := git_mock.NewMockGitService(ctrl)

			// Setup mocks
			tt.setupMocks(mockHosting, mockWorkspace, mockLayout, mockOverlay, mockRefParser, mockGit)

			// Create UseCase to test
			useCase := testtarget.NewUseCase(mockHosting, mockWorkspace, mockOverlay, mockRefParser, mockGit)

			// Execute test
			err := useCase.Execute(context.Background(), tt.refWithAlias, tt.template, tt.options)

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
