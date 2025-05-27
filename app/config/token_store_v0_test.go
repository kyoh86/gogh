package config_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

// setupTokenStoreV0Test creates a temporary directory for testing and returns the necessary test components.
func setupTokenStoreV0Test(t *testing.T) (
	string,
	func(),
	*auth_mock.MockTokenService,
	*config.TokenStoreV0,
) {
	ctrl := gomock.NewController(t)
	tempDir, err := os.MkdirTemp("", "token-store-v0-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	mockService := auth_mock.NewMockTokenService(ctrl)
	store := config.NewTokenStoreV0()

	// Override the appContextPath to use our test directory
	origAppContextPath := config.AppContextPathFunc
	config.AppContextPathFunc = func(envName string, fallbackFunc func() (string, error), rel ...string) (string, error) {
		return filepath.Join(append([]string{tempDir}, rel...)...), nil
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
		// Restore the original AppContextPathFunc
		config.AppContextPathFunc = origAppContextPath
	}

	return tempDir, cleanup, mockService, store
}

// createTestTokenYAML creates a test YAML file with token data.
func createTestTokenYAML(t *testing.T, path string) {
	t.Helper()

	content := `hosts:
  github.com:
    owners:
      testuser:
        access_token: test-access-token
        token_type: Bearer
        refresh_token: test-refresh-token
        expiry: 2023-12-31T23:59:59Z
      otheruser:
        access_token: other-access-token
        token_type: Bearer
  gitlab.com:
    owners:
      testuser:
        access_token: gitlab-access-token
`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test YAML file: %v", err)
	}
}

func TestNewTokenStoreV0(t *testing.T) {
	store := config.NewTokenStoreV0()
	if store == nil {
		t.Fatal("Expected non-nil store")
	}
}

func TestTokenStoreV0Source(t *testing.T) {
	tempDir, cleanup, _, store := setupTokenStoreV0Test(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "tokens.yaml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestTokenStoreV0Load_FileExists(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreV0Test(t)
	defer cleanup()

	// Create a test YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	createTestTokenYAML(t, configPath)

	// Setup mock expectations
	mockService.EXPECT().Set("github.com", "testuser", gomock.Any()).DoAndReturn(
		func(host, owner string, token oauth2.Token) error {
			if token.AccessToken != "test-access-token" {
				t.Errorf("Expected access token %q, got %q", "test-access-token", token.AccessToken)
			}
			if token.TokenType != "Bearer" {
				t.Errorf("Expected token type %q, got %q", "Bearer", token.TokenType)
			}
			if token.RefreshToken != "test-refresh-token" {
				t.Errorf("Expected refresh token %q, got %q", "test-refresh-token", token.RefreshToken)
			}
			return nil
		})

	mockService.EXPECT().Set("github.com", "otheruser", gomock.Any()).DoAndReturn(
		func(host, owner string, token oauth2.Token) error {
			if token.AccessToken != "other-access-token" {
				t.Errorf("Expected access token %q, got %q", "other-access-token", token.AccessToken)
			}
			return nil
		})

	mockService.EXPECT().Set("gitlab.com", "testuser", gomock.Any()).DoAndReturn(
		func(host, owner string, token oauth2.Token) error {
			if token.AccessToken != "gitlab-access-token" {
				t.Errorf("Expected access token %q, got %q", "gitlab-access-token", token.AccessToken)
			}
			return nil
		})

	mockService.EXPECT().MarkSaved()

	// Call Load
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	service, err := store.Load(ctx, initialFunc)
	if err != nil {
		t.Fatalf("Unexpected error from Load(): %v", err)
	}

	if service != mockService {
		t.Error("Expected Load to return the service from initialFunc")
	}
}

func TestTokenStoreV0Load_FileDoesNotExist(t *testing.T) {
	_, cleanup, mockService, store := setupTokenStoreV0Test(t)
	defer cleanup()

	// Call Load with no file
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() when file doesn't exist, got nil")
	}
}

func TestTokenStoreV0Load_EmptyAccessToken(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreV0Test(t)
	defer cleanup()

	// Create a YAML file with an empty access token
	configPath := filepath.Join(tempDir, "tokens.yaml")
	content := `hosts:
  github.com:
    owners:
      testuser:
        access_token: ""
        token_type: Bearer
      validuser:
        access_token: valid-token
        token_type: Bearer
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	// Setup mock expectations - should skip the empty token
	mockService.EXPECT().Set("github.com", "validuser", gomock.Any()).Return(nil)
	mockService.EXPECT().MarkSaved()

	// Call Load
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err != nil {
		t.Fatalf("Unexpected error from Load(): %v", err)
	}
}

func TestTokenStoreV0Load_SetError(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreV0Test(t)
	defer cleanup()

	// Create a test YAML file
	configPath := filepath.Join(tempDir, "tokens.yaml")
	createTestTokenYAML(t, configPath)

	// Setup mock expectations - Set call will return an error
	expectedErr := errors.New("test error")
	mockService.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)

	// Call Load
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	_, err := store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() when Set fails, got nil")
	}
}
