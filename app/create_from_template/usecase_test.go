package create_from_template_test

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/create_from_template"
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
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	// Define test cases
	tests := []struct {
		name         string
		refWithAlias string
		template     repository.Reference
		options      testtarget.CreateFromTemplateOptions
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
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
				mockRefParser *repository_mock.MockReferenceParser,
				mockGit *git_mock.MockGitService,
			) {
				tmpDir := t.TempDir()
				ref := repository.NewReference("github.com", "kyoh86", "new-repo")
				templateRef := repository.NewReference("github.com", "kyoh86", "template-repo")

				// Parse reference
				mockRefParser.EXPECT().
					ParseWithAlias("github.com/kyoh86/new-repo").
					Return(&repository.ReferenceWithAlias{
						Reference: ref,
						Alias:     nil,
					}, nil).AnyTimes()

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

				pseudoPath := filepath.Join(tmpDir, "github.com/kyoh86/gogh")
				// Get primary layout to clone into
				mockLayout.EXPECT().PathFor(ref).Return(pseudoPath)
				mockWorkspace.EXPECT().GetPrimaryLayout().Return(mockLayout)

				// Set expectations for GitService.Clone
				mockGit.EXPECT().AuthenticateWithUsernamePassword(gomock.Any(), "kyoh86", "").Return(mockGit, nil)
				mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				mockHook.EXPECT().
					ListFor(ref, hook.EventPostCreate).Return(func(yield func(hook.Hook, error) bool) {
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-create-from-template-example",
						OperationType: hook.OperationTypeOverlay,
						OperationID:   "overlay-id",
					}), nil) {
						return
					}
					if !yield(hook.NewHook(hook.Entry{
						Name:          "post-create-from-template-example",
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
			errorContains: "script error",
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
				mockFinder *workspace_mock.MockFinderService,
				mockLayout *workspace_mock.MockLayoutService,
				mockOverlay *overlay_mock.MockOverlayService,
				mockScript *script_mock.MockScriptService,
				mockHook *hook_mock.MockHookService,
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
			usecase := testtarget.NewUsecase(mockHosting, mockWorkspace, mockFinder, mockOverlay, mockScript, mockHook, mockRefParser, mockGit)

			// Execute test
			err := usecase.Execute(context.Background(), tt.refWithAlias, tt.template, tt.options)

			// Verify results
			if err == nil {
				t.Errorf("Expected an error but got none")
			} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("Expected error containing: %v, got: %v", tt.errorContains, err)
			}
		})
	}
}
