package repos_test

import (
	"context"
	"testing"

	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	t.Run("successfully list repositories", func(t *testing.T) {
		// Setup mock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHostingService := hosting_mock.NewMockHostingService(ctrl)
		expectedOpts := hosting.ListRepositoryOptions{
			Limit: 30, // Default limit
		}

		// Mock repositories to return
		repositories := []*hosting.Repository{
			{Ref: repository.NewReference("github.com", "kyoh86", "repo1")},
			{Ref: repository.NewReference("github.com", "kyoh86", "repo2")},
		}

		// Set expectations
		mockHostingService.EXPECT().
			ListRepository(gomock.Any(), expectedOpts).
			Return(func(yield func(*hosting.Repository, error) bool) {
				for _, repo := range repositories {
					if !yield(repo, nil) {
						return
					}
				}
			})

		// Create Usecase with mock
		usecase := repos.NewUsecase(mockHostingService)

		// Execute the method
		options := repos.Options{} // Default options
		results := []*hosting.Repository{}
		var testErr error

		for repo, err := range usecase.Execute(context.Background(), options) {
			if err != nil {
				testErr = err
				break
			}
			results = append(results, repo)
		}

		// Verify results
		if testErr != nil {
			t.Fatalf("Expected no error, got %v", testErr)
		}

		if len(results) != len(repositories) {
			t.Fatalf("Expected %d repositories, got %d", len(repositories), len(results))
		}

		for i, repo := range results {
			if repo.Ref != repositories[i].Ref {
				t.Errorf("Expected repo %s, got %s", repositories[i].Ref, repo.Ref)
			}
		}
	})

	t.Run("error in options conversion", func(t *testing.T) {
		// Setup mock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHostingService := hosting_mock.NewMockHostingService(ctrl)
		usecase := repos.NewUsecase(mockHostingService)

		// Invalid options that will cause an error
		options := repos.Options{
			Privacy: "invalid-privacy", // This should cause an error
		}

		var testErr error
		count := 0

		for _, err := range usecase.Execute(context.Background(), options) {
			count++
			if err != nil {
				testErr = err
			}
		}

		// Should get an error and stop after first yield
		if testErr == nil {
			t.Fatalf("Expected an error, got none")
		}
		if count != 1 {
			t.Fatalf("Expected 1 iteration, got %d", count)
		}
	})
}
