package github_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	testtarget "github.com/kyoh86/gogh/v4/infra/github"
	"go.uber.org/mock/gomock"
)

// These tests are limited because we can't easily mock the GitHub API client
// In a real test, you'd use a custom mock or a test double for the GitHub client

// Setup test data and mocks
func setupHostingServiceTest(t *testing.T) (*gomock.Controller, *auth_mock.MockTokenService, *repository_mock.MockDefaultNameService, *testtarget.HostingService) {
	ctrl := gomock.NewController(t)
	mockTokenService := auth_mock.NewMockTokenService(ctrl)
	mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
	service := testtarget.NewHostingService(mockTokenService, mockDefaultNameService)

	return ctrl, mockTokenService, mockDefaultNameService, service
}

func TestGetURLOf(t *testing.T) {
	ctrl, _, _, service := setupHostingServiceTest(t)
	defer ctrl.Finish()

	testCases := []struct {
		name        string
		ref         repository.Reference
		expectedURL string
	}{
		{
			name:        "github repository",
			ref:         repository.NewReference("github.com", "kyoh86", "gogh"),
			expectedURL: "https://github.com/kyoh86/gogh",
		},
		{
			name:        "enterprise github",
			ref:         repository.NewReference("github.mycompany.com", "user", "project"),
			expectedURL: "https://github.mycompany.com/user/project",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := service.GetURLOf(tc.ref)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if u.String() != tc.expectedURL {
				t.Errorf("Expected URL %q, got %q", tc.expectedURL, u.String())
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	ctrl, _, _, service := setupHostingServiceTest(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		urlStr        string
		expectedHost  string
		expectedOwner string
		expectedName  string
		expectError   bool
	}{
		{
			name:          "valid github URL",
			urlStr:        "https://github.com/kyoh86/gogh",
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
			expectError:   false,
		},
		{
			name:          "valid URL with .git suffix",
			urlStr:        "https://github.com/kyoh86/gogh.git",
			expectedHost:  "github.com",
			expectedOwner: "kyoh86",
			expectedName:  "gogh",
			expectError:   false,
		},
		{
			name:          "valid enterprise github URL",
			urlStr:        "https://github.mycompany.com/user/project",
			expectedHost:  "github.mycompany.com",
			expectedOwner: "user",
			expectedName:  "project",
			expectError:   false,
		},
		{
			name:        "invalid URL with missing path components",
			urlStr:      "https://github.com/kyoh86",
			expectError: true,
		},
		{
			name:        "invalid URL with empty path",
			urlStr:      "https://github.com",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, _ := url.Parse(tc.urlStr)
			ref, err := service.ParseURL(u)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if ref.Host() != tc.expectedHost {
				t.Errorf("Expected host %q, got %q", tc.expectedHost, ref.Host())
			}

			if ref.Owner() != tc.expectedOwner {
				t.Errorf("Expected owner %q, got %q", tc.expectedOwner, ref.Owner())
			}

			if ref.Name() != tc.expectedName {
				t.Errorf("Expected name %q, got %q", tc.expectedName, ref.Name())
			}
		})
	}
}

func TestGetTokenFor(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Token exists for exact host/owner
	t.Run("exact match", func(t *testing.T) {
		ctrl, mockTokenService, _, service := setupHostingServiceTest(t)
		defer ctrl.Finish()

		ref := repository.NewReference("github.com", "user1", "repo")
		token := auth.Token{AccessToken: "token1"}

		mockTokenService.EXPECT().Has("github.com", "user1").Return(true)
		mockTokenService.EXPECT().Get("github.com", "user1").Return(token, nil)

		tokenOwner, gotToken, err := service.GetTokenFor(ctx, ref.Host(), ref.Owner())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if tokenOwner != "user1" {
			t.Errorf("Expected token owner %q, got %q", "user1", tokenOwner)
		}

		if gotToken != token {
			t.Errorf("Expected token %v, got %v", token, gotToken)
		}
	})

	// Test case 2: No token available
	t.Run("no token", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()

		ref := repository.NewReference("github.com", "unknown", "repo")

		mockTokenService.EXPECT().Has("github.com", "unknown").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		_, _, err := service.GetTokenFor(ctx, ref.Host(), ref.Owner())
		if err == nil {
			t.Error("Expected error for missing token, got nil")
		}
	})
}

func TestGetRepository(t *testing.T) {
	ctx := context.Background()

	// Test error handling for token retrieval
	t.Run("token error", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()

		// Mock token lookup
		mockTokenService.EXPECT().Has("github.com", "error").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		errorRef := repository.NewReference("github.com", "error", "repo")
		_, err := service.GetRepository(ctx, errorRef)
		if err == nil {
			t.Error("Expected error for token retrieval, got nil")
		}
	})
}

func TestListRepository(t *testing.T) {
	ctrl, mockTokenService, _, service := setupHostingServiceTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	token := auth.Token{AccessToken: "test-token"}

	// Setup token entries
	entries := []auth.TokenEntry{
		{Host: "github.com", Owner: "user1", Token: token},
	}

	mockTokenService.EXPECT().Entries().Return(entries).AnyTimes()

	// This test is limited because we can't easily mock the GraphQL API client
	// In a real test, you'd use a custom mock for the GraphQL client

	// Basic test for the list repository function
	opts := hosting.ListRepositoryOptions{
		Limit: 10,
	}

	// Just test that the function doesn't panic
	count := 0
	for _, err := range service.ListRepository(ctx, opts) {
		if err != nil {
			// We expect errors without mocking the GraphQL client
			continue
		}
		count++
	}
}

func TestCreateRepository(t *testing.T) {
	ctx := context.Background()

	// Test error handling for token retrieval
	t.Run("token error", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()
		mockTokenService.EXPECT().Has("github.com", "error").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		errorRef := repository.NewReference("github.com", "error", "repo")
		_, err := service.CreateRepository(ctx, errorRef, hosting.CreateRepositoryOptions{})
		if err == nil {
			t.Error("Expected error for token retrieval, got nil")
		}
	})
}

func TestCreateRepositoryFromTemplate(t *testing.T) {
	ctx := context.Background()
	template := repository.NewReference("github.com", "kyoh86", "template-repo")

	// Test error handling for token retrieval
	t.Run("token error", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()
		mockTokenService.EXPECT().Has("github.com", "error").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		errorRef := repository.NewReference("github.com", "error", "repo")
		_, err := service.CreateRepositoryFromTemplate(ctx, errorRef, template, hosting.CreateRepositoryFromTemplateOptions{})
		if err == nil {
			t.Error("Expected error for token retrieval, got nil")
		}
	})
}

func TestDeleteRepository(t *testing.T) {
	ctx := context.Background()

	// Test error handling for token retrieval
	t.Run("token error", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()
		mockTokenService.EXPECT().Has("github.com", "error").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		errorRef := repository.NewReference("github.com", "error", "repo")
		err := service.DeleteRepository(ctx, errorRef)
		if err == nil {
			t.Error("Expected error for token retrieval, got nil")
		}
	})
}

func TestForkRepository(t *testing.T) {
	ctx := context.Background()
	sourceRef := repository.NewReference("github.com", "original", "repo")

	// Test error handling for token retrieval
	t.Run("token error", func(t *testing.T) {
		ctrl, mockTokenService, mockDefaultNameService, service := setupHostingServiceTest(t)
		defer ctrl.Finish()
		mockTokenService.EXPECT().Has("github.com", "error").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		mockDefaultNameService.EXPECT().GetDefaultOwnerFor("github.com").Return("default-owner", nil)

		mockTokenService.EXPECT().Has("github.com", "default-owner").Return(false)
		mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})

		errorRef := repository.NewReference("github.com", "error", "repo")
		_, err := service.ForkRepository(ctx, sourceRef, errorRef, hosting.ForkRepositoryOptions{})
		if err == nil {
			t.Error("Expected error for token retrieval, got nil")
		}
	})
}
