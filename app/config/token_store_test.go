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

// setupTokenStoreTest creates a temporary directory for testing and returns the necessary test components.
func setupTokenStoreTest(t *testing.T) (
	string,
	func(),
	*auth_mock.MockTokenService,
	*config.TokenStore,
) {
	ctrl := gomock.NewController(t)
	tempDir, err := os.MkdirTemp("", "token-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	mockService := auth_mock.NewMockTokenService(ctrl)
	store := config.NewTokenStore()

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

// createTestTokenTOML creates a test TOML file with token data.
func createTestTokenTOML(t *testing.T, path string) {
	t.Helper()

	content := `['github.com'.testuser]
AccessToken = "test-access-token"
TokenType = "Bearer"
RefreshToken = "test-refresh-token"
Expiry = 2023-12-31T23:59:59Z

['github.com'.otheruser]
AccessToken = "other-access-token"
TokenType = "Bearer"

['gitlab.com'.testuser]
AccessToken = "gitlab-access-token"
`

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test TOML file: %v", err)
	}
}

func TestTokenStoreSource(t *testing.T) {
	tempDir, cleanup, _, store := setupTokenStoreTest(t)
	defer cleanup()

	path, err := store.Source()
	if err != nil {
		t.Fatalf("Unexpected error from Source(): %v", err)
	}

	expectedPath := filepath.Join(tempDir, "tokens.v4.toml")
	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestTokenStoreLoad_FileExists(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	createTestTokenTOML(t, configPath)

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

func TestTokenStoreLoad_FileDoesNotExist(t *testing.T) {
	_, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Call Load with no file
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	if _, err := store.Load(ctx, initialFunc); err == nil {
		t.Errorf("Expected error from Load() when file doesn't exist: %v", err)
	}
}

func TestTokenStoreLoad_InvalidTOML(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Create an invalid TOML file
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	err := os.WriteFile(configPath, []byte("invalid toml content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid TOML file: %v", err)
	}

	// Call Load
	ctx := context.Background()
	initialFunc := func() auth.TokenService {
		return mockService
	}

	_, err = store.Load(ctx, initialFunc)
	if err == nil {
		t.Fatal("Expected error from Load() with invalid TOML, got nil")
	}
}

func TestTokenStoreLoad_EmptyAccessToken(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Create a TOML file with an empty access token
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	content := `['github.com'.testuser]
AccessToken = ""
TokenType = "Bearer"

['github.com'.validuser]
AccessToken = "valid-token"
TokenType = "Bearer"
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write TOML file: %v", err)
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

func TestTokenStoreLoad_SetError(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Create a test TOML file
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	createTestTokenTOML(t, configPath)

	// Setup mock expectations - second Set call will return an error
	expectedErr := errors.New("test error")
	mockService.EXPECT().Set("github.com", "testuser", gomock.Any()).Return(expectedErr).AnyTimes()
	mockService.EXPECT().Set("github.com", "otheruser", gomock.Any()).Return(expectedErr).AnyTimes()
	mockService.EXPECT().Set("gitlab.com", "testuser", gomock.Any()).Return(expectedErr).AnyTimes()

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

func TestTokenStoreSave_NoChanges(t *testing.T) {
	_, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}
}

func TestTokenStoreSave_WithChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Setup mock token entries
	token1 := oauth2.Token{
		AccessToken:  "access1",
		TokenType:    "Bearer",
		RefreshToken: "refresh1",
	}
	token2 := oauth2.Token{
		AccessToken: "access2",
	}

	entries := []auth.TokenEntry{
		{Host: "github.com", Owner: "user1", Token: token1},
		{Host: "gitlab.com", Owner: "user2", Token: token2},
	}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(true)
	mockService.EXPECT().Entries().Return(entries)
	mockService.EXPECT().MarkSaved()

	// Call Save
	ctx := context.Background()
	err := store.Save(ctx, mockService, false)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}

func TestTokenStoreSave_ForceWithoutChanges(t *testing.T) {
	tempDir, cleanup, mockService, store := setupTokenStoreTest(t)
	defer cleanup()

	// Setup mock token entries
	token := oauth2.Token{AccessToken: "test-token"}
	entries := []auth.TokenEntry{
		{Host: "github.com", Owner: "testuser", Token: token},
	}

	// Setup mock expectations
	mockService.EXPECT().HasChanges().Return(false)
	mockService.EXPECT().Entries().Return(entries)
	mockService.EXPECT().MarkSaved()

	// Call Save with force=true
	ctx := context.Background()
	err := store.Save(ctx, mockService, true)
	if err != nil {
		t.Fatalf("Unexpected error from Save(): %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tempDir, "tokens.v4.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved TOML file: %v", err)
	}

	// Check that the content looks reasonable
	if len(content) == 0 {
		t.Error("Saved TOML file is empty")
	}
}
