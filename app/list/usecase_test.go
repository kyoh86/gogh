package list_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"iter"

	testtarget "github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		options       testtarget.Options
		setupMocks    func(mockWorkspace *workspace_mock.MockWorkspaceService, mockFinder *workspace_mock.MockFinderService)
		expectedError bool
		errorContains string
	}{
		{
			name: "List repositories from primary root",
			options: testtarget.Options{
				Primary:     true,
				ListOptions: workspace.ListOptions{},
			},
			setupMocks: func(mockWorkspace *workspace_mock.MockWorkspaceService, mockFinder *workspace_mock.MockFinderService) {
				primaryRoot := "/path/to/primary"
				mockLayout := workspace_mock.NewMockLayoutService(gomock.NewController(t))

				// Set up expectations for primary root listing
				mockWorkspace.EXPECT().GetPrimaryRoot().Return(primaryRoot)
				mockWorkspace.EXPECT().GetLayoutFor(primaryRoot).Return(mockLayout)

				// Create a sample repository sequence
				repoSeq := func(yield func(*repository.Location, error) bool) {
					yield(repository.NewLocation(
						"/path/to/primary/github.com/kyoh86/repo1",
						"github.com", "kyoh86", "repo1",
					), nil)
					yield(repository.NewLocation(
						"/path/to/primary/github.com/kyoh86/repo2",
						"github.com", "kyoh86", "repo2",
					), nil)
				}

				// Expect ListRepositoryInRoot to be called
				mockFinder.EXPECT().
					ListRepositoryInRoot(gomock.Any(), mockLayout, workspace.ListOptions{}).
					Return(iter.Seq2[*repository.Location, error](repoSeq))
			},
			expectedError: false,
		},
		{
			name: "List all repositories",
			options: testtarget.Options{
				Primary:     false,
				ListOptions: workspace.ListOptions{},
			},
			setupMocks: func(mockWorkspace *workspace_mock.MockWorkspaceService, mockFinder *workspace_mock.MockFinderService) {
				// Create a sample repository sequence
				repoSeq := func(yield func(*repository.Location, error) bool) {
					yield(repository.NewLocation(
						"/path/to/root1/github.com/kyoh86/repo1",
						"github.com", "kyoh86", "repo1",
					), nil)
					yield(repository.NewLocation(
						"/path/to/root2/github.com/kyoh86/repo2",
						"github.com", "kyoh86", "repo2",
					), nil)
				}

				// Expect ListAllRepository to be called
				mockFinder.EXPECT().
					ListAllRepository(gomock.Any(), mockWorkspace, workspace.ListOptions{}).
					Return(iter.Seq2[*repository.Location, error](repoSeq))
			},
			expectedError: false,
		},
		{
			name: "Error during repository listing",
			options: testtarget.Options{
				Primary:     false,
				ListOptions: workspace.ListOptions{},
			},
			setupMocks: func(mockWorkspace *workspace_mock.MockWorkspaceService, mockFinder *workspace_mock.MockFinderService) {
				// Create an error sequence
				errorSeq := func(yield func(*repository.Location, error) bool) {
					yield(nil, errors.New("failed to list repositories"))
				}

				// Expect ListAllRepository to be called and return an error
				mockFinder.EXPECT().
					ListAllRepository(gomock.Any(), mockWorkspace, workspace.ListOptions{}).
					Return(iter.Seq2[*repository.Location, error](errorSeq))
			},
			expectedError: true,
			errorContains: "failed to list repositories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
			mockWorkspace := workspace_mock.NewMockWorkspaceService(ctrl)
			mockFinder := workspace_mock.NewMockFinderService(ctrl)

			// Setup mocks
			tt.setupMocks(mockWorkspace, mockFinder)

			// Create Usecase to test
			usecase := testtarget.NewUsecase(mockWorkspace, mockFinder)

			// Execute test
			seq := usecase.Execute(context.Background(), tt.options)

			// Process the sequence to check for errors
			var foundError error
			for item, err := range seq {
				if err != nil {
					foundError = err
					break
				}
				if item == nil {
					t.Errorf("Expected repository location, got nil")
				}
			}

			// Verify results
			if tt.expectedError {
				if foundError == nil {
					t.Errorf("Expected an error but got none")
				} else if tt.errorContains != "" && !strings.Contains(foundError.Error(), tt.errorContains) {
					t.Errorf("Expected error containing: %v, got: %v", tt.errorContains, foundError)
				}
			} else {
				if foundError != nil {
					t.Errorf("Unexpected error: %v", foundError)
				}
			}
		})
	}
}
