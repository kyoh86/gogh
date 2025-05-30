package bundle_dump_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/bundle_dump"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestExecute(t *testing.T) {
	// Test cases
	tests := []struct {
		name            string
		setupMocks      func(*workspace_mock.MockFinderService, *workspace_mock.MockWorkspaceService, *hosting_mock.MockHostingService, *git_mock.MockGitService)
		options         workspace.ListOptions
		expectedCount   int
		expectedError   bool
		expectedEntries []*testtarget.BundleEntry
	}{
		{
			name: "Success: When repository can be retrieved",
			setupMocks: func(finder *workspace_mock.MockFinderService, ws *workspace_mock.MockWorkspaceService, hosting *hosting_mock.MockHostingService, git *git_mock.MockGitService) {
				// Setup repository
				loc := repository.NewLocation("/path/to/github.com/kyoh86/gogh", "github.com", "kyoh86", "gogh")

				// Setup mock for ListAllRepository
				finder.EXPECT().
					ListAllRepository(gomock.Any(), ws, gomock.Any()).
					Return(func(yield func(*repository.Location, error) bool) {
						yield(loc, nil)
					})

				// Setup mock for GetDefaultRemotes
				git.EXPECT().
					GetDefaultRemotes(gomock.Any(), "/path/to/github.com/kyoh86/gogh").
					Return([]string{"https://github.com/kyoh86/gogh.git"}, nil)

				// Setup mock for ParseURL
				expectedURL, _ := url.Parse("https://github.com/kyoh86/gogh.git")
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				hosting.EXPECT().
					ParseURL(gomock.Eq(expectedURL)).
					Return(&ref, nil)
			},
			options:       workspace.ListOptions{},
			expectedCount: 1,
			expectedError: false,
			expectedEntries: []*testtarget.BundleEntry{
				{
					Name:  "github.com/kyoh86/gogh",
					Alias: nil,
				},
			},
		},
		{
			name: "Success: When remote name differs from local path",
			setupMocks: func(finder *workspace_mock.MockFinderService, ws *workspace_mock.MockWorkspaceService, hosting *hosting_mock.MockHostingService, git *git_mock.MockGitService) {
				// Setup repository
				loc := repository.NewLocation("/path/to/github.com/user/fork-repo", "github.com", "user", "fork-repo")

				// Setup mock for ListAllRepository
				finder.EXPECT().
					ListAllRepository(gomock.Any(), ws, gomock.Any()).
					Return(func(yield func(*repository.Location, error) bool) {
						yield(loc, nil)
					})

				// Setup mock for GetDefaultRemotes
				git.EXPECT().
					GetDefaultRemotes(gomock.Any(), "/path/to/github.com/user/fork-repo").
					Return([]string{"https://github.com/original/repo.git"}, nil)

				// Setup mock for ParseURL
				expectedURL, _ := url.Parse("https://github.com/original/repo.git")
				ref := repository.NewReference("github.com", "original", "repo")
				hosting.EXPECT().
					ParseURL(gomock.Eq(expectedURL)).
					Return(&ref, nil)
			},
			options:       workspace.ListOptions{},
			expectedCount: 1,
			expectedError: false,
			expectedEntries: []*testtarget.BundleEntry{
				{
					Name:  "github.com/user/fork-repo",
					Alias: stringPtr("github.com/original/repo"),
				},
			},
		},
		{
			name: "Error: When FinderService returns an error",
			setupMocks: func(finder *workspace_mock.MockFinderService, ws *workspace_mock.MockWorkspaceService, hosting *hosting_mock.MockHostingService, git *git_mock.MockGitService) {
				// Setup mock for ListAllRepository to return an error
				finder.EXPECT().
					ListAllRepository(gomock.Any(), ws, gomock.Any()).
					Return(func(yield func(*repository.Location, error) bool) {
						yield(nil, errors.New("repository find error"))
					})
			},
			options:       workspace.ListOptions{},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mocks
			mockFinder := workspace_mock.NewMockFinderService(ctrl)
			mockWs := workspace_mock.NewMockWorkspaceService(ctrl)
			mockHosting := hosting_mock.NewMockHostingService(ctrl)
			mockGit := git_mock.NewMockGitService(ctrl)

			// Setup mocks
			tt.setupMocks(mockFinder, mockWs, mockHosting, mockGit)

			// Create the target UseCase
			useCase := testtarget.NewUseCase(mockWs, mockFinder, mockHosting, mockGit)

			// Execute the UseCase
			ctx := context.Background()
			result := useCase.Execute(ctx, tt.options)

			// Verify the results
			count := 0
			var err error
			for entry, iterErr := range result {
				if iterErr != nil {
					err = iterErr
					break
				}
				count++

				// Compare expected and actual values
				if !tt.expectedError && count <= len(tt.expectedEntries) {
					expected := tt.expectedEntries[count-1]
					if expected.Name != entry.Name {
						t.Errorf("Name mismatch: expected %s, got %s", expected.Name, entry.Name)
					}

					if (expected.Alias == nil && entry.Alias != nil) ||
						(expected.Alias != nil && entry.Alias == nil) ||
						(expected.Alias != nil && entry.Alias != nil && *expected.Alias != *entry.Alias) {
						t.Errorf("Alias mismatch: expected %v, got %v",
							expected.Alias, entry.Alias)
					}
				}
			}

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectedError && count != tt.expectedCount {
				t.Errorf("Expected %d entries, got %d", tt.expectedCount, count)
			}
		})
	}
}

// Helper function to return a pointer to a string
func stringPtr(s string) *string {
	return &s
}
