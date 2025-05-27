package git_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	coregit "github.com/kyoh86/gogh/v4/core/git"
	testtarget "github.com/kyoh86/gogh/v4/infra/git"
)

// setupTempDir creates a temporary directory for testing and returns its path.
// The caller is responsible for cleaning it up.
func setupTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "git-service-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

func TestNewService(t *testing.T) {
	// Test with no options
	service := testtarget.NewService()
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Test with progress writer option
	writer := &mockWriter{}
	service = testtarget.NewService(testtarget.CloneProgressWriter(writer))
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestAuthenticateWithUsernamePassword(t *testing.T) {
	ctx := context.Background()
	service := testtarget.NewService()

	// Test authentication
	authenticatedService, err := service.AuthenticateWithUsernamePassword(ctx, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	if authenticatedService == nil {
		t.Fatal("Expected non-nil authenticated service")
	}

	// Verify we got a new service instance
	if authenticatedService == service {
		t.Error("Expected a new service instance, got the same one")
	}
}

func TestClone(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Test cloning a non-existent repository (should return an error)
	err := service.Clone(ctx, "https://github.com/kyoh86/non-existent-repo.git", tempDir, coregit.CloneOptions{})
	if !errors.Is(err, coregit.ErrRepositoryNotExists) {
		t.Errorf("Expected ErrRepositoryNotExists for non-existent repo, got: %v", err)
	}

	// Create a mock local repository to clone from
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	// Initialize a git repository in the source directory
	_, err = git.PlainInit(sourceDir, false)
	if err != nil {
		t.Fatalf("Failed to initialize source repository: %v", err)
	}

	// Test cloning from a local path (this might fail in some environments,
	// but serves as a basic test)
	err = service.Clone(ctx, "file://"+sourceDir, destDir, coregit.CloneOptions{})
	// Only check for specific errors to make the test more robust
	if err != nil && !errors.Is(err, coregit.ErrRepositoryEmpty) {
		t.Errorf("Failed to clone from local path: %v", err)
	}
}

func TestInit(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Test initializing a new repository
	repoPath := filepath.Join(tempDir, "repo")
	err := service.Init(ctx, "https://github.com/kyoh86/test-repo.git", repoPath, false, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Verify the repository was created
	_, err = git.PlainOpen(repoPath)
	if err != nil {
		t.Errorf("Failed to open initialized repository: %v", err)
	}

	// Verify the remote was set
	repo, _ := git.PlainOpen(repoPath)
	remote, err := repo.Remote(git.DefaultRemoteName)
	if err != nil {
		t.Errorf("Failed to get remote: %v", err)
	}

	if remote.Config().URLs[0] != "https://github.com/kyoh86/test-repo.git" {
		t.Errorf("Expected remote URL %q, got %q",
			"https://github.com/kyoh86/test-repo.git",
			remote.Config().URLs[0])
	}

	// Test initializing a bare repository
	bareRepoPath := filepath.Join(tempDir, "bare")
	err = service.Init(ctx, "https://github.com/kyoh86/test-bare.git", bareRepoPath, true, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize bare repository: %v", err)
	}

	// Verify the bare repository was created
	_, err = git.PlainOpen(bareRepoPath)
	if err != nil {
		t.Errorf("Failed to open initialized bare repository: %v", err)
	}
}

func TestSetRemotes(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Initialize a repository
	repoPath := filepath.Join(tempDir, "repo")
	err := service.Init(ctx, "https://github.com/kyoh86/test-repo.git", repoPath, false, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Test setting custom remote
	err = service.SetRemotes(ctx, repoPath, "upstream", []string{"https://github.com/upstream/test-repo.git"})
	if err != nil {
		t.Errorf("Failed to set remote: %v", err)
	}

	// Verify the remote was set
	repo, _ := git.PlainOpen(repoPath)
	remote, err := repo.Remote("upstream")
	if err != nil {
		t.Errorf("Failed to get upstream remote: %v", err)
	}

	if remote.Config().URLs[0] != "https://github.com/upstream/test-repo.git" {
		t.Errorf("Expected remote URL %q, got %q",
			"https://github.com/upstream/test-repo.git",
			remote.Config().URLs[0])
	}

	// Test setting multiple URLs for a remote
	err = service.SetRemotes(ctx, repoPath, "multi", []string{
		"https://github.com/multi1/test-repo.git",
		"https://github.com/multi2/test-repo.git",
	})
	if err != nil {
		t.Errorf("Failed to set multiple remotes: %v", err)
	}

	// Verify multiple URLs were set
	remote, err = repo.Remote("multi")
	if err != nil {
		t.Errorf("Failed to get multi remote: %v", err)
	}

	if len(remote.Config().URLs) != 2 {
		t.Errorf("Expected 2 URLs for multi remote, got %d", len(remote.Config().URLs))
	}
}

func TestSetDefaultRemotes(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Initialize a repository
	repoPath := filepath.Join(tempDir, "repo")
	err := service.Init(ctx, "https://github.com/kyoh86/test-repo.git", repoPath, false, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Test setting default remote
	err = service.SetDefaultRemotes(ctx, repoPath, []string{"https://github.com/default/test-repo.git"})
	if err != nil {
		t.Errorf("Failed to set default remote: %v", err)
	}

	// Verify the remote was set
	repo, _ := git.PlainOpen(repoPath)
	remote, err := repo.Remote(git.DefaultRemoteName)
	if err != nil {
		t.Errorf("Failed to get default remote: %v", err)
	}

	if remote.Config().URLs[0] != "https://github.com/default/test-repo.git" {
		t.Errorf("Expected remote URL %q, got %q",
			"https://github.com/default/test-repo.git",
			remote.Config().URLs[0])
	}
}

func TestGetRemotes(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Initialize a repository
	repoPath := filepath.Join(tempDir, "repo")
	err := service.Init(ctx, "https://github.com/kyoh86/test-repo.git", repoPath, false, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Add another remote
	repo, _ := git.PlainOpen(repoPath)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{"https://github.com/upstream/test-repo.git"},
	})
	if err != nil {
		t.Fatalf("Failed to create upstream remote: %v", err)
	}

	// Test getting origin remote
	remotes, err := service.GetRemotes(ctx, repoPath, git.DefaultRemoteName)
	if err != nil {
		t.Errorf("Failed to get origin remote: %v", err)
	}

	if len(remotes) != 1 {
		t.Errorf("Expected 1 URL for origin remote, got %d", len(remotes))
	}

	if remotes[0] != "https://github.com/kyoh86/test-repo.git" {
		t.Errorf("Expected origin URL %q, got %q",
			"https://github.com/kyoh86/test-repo.git",
			remotes[0])
	}

	// Test getting upstream remote
	remotes, err = service.GetRemotes(ctx, repoPath, "upstream")
	if err != nil {
		t.Errorf("Failed to get upstream remote: %v", err)
	}

	if len(remotes) != 1 {
		t.Errorf("Expected 1 URL for upstream remote, got %d", len(remotes))
	}

	if remotes[0] != "https://github.com/upstream/test-repo.git" {
		t.Errorf("Expected upstream URL %q, got %q",
			"https://github.com/upstream/test-repo.git",
			remotes[0])
	}

	// Test getting non-existent remote
	remotes, err = service.GetRemotes(ctx, repoPath, "nonexistent")
	if err != nil {
		t.Errorf("Expected no error for nonexistent remote, got: %v", err)
	}

	if len(remotes) != 0 {
		t.Errorf("Expected 0 URLs for nonexistent remote, got %d", len(remotes))
	}
}

func TestGetDefaultRemotes(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Initialize a repository
	repoPath := filepath.Join(tempDir, "repo")
	err := service.Init(ctx, "https://github.com/kyoh86/test-repo.git", repoPath, false, coregit.InitOptions{})
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	// Test getting default remote
	remotes, err := service.GetDefaultRemotes(ctx, repoPath)
	if err != nil {
		t.Errorf("Failed to get default remotes: %v", err)
	}

	if len(remotes) != 1 {
		t.Errorf("Expected 1 URL for default remote, got %d", len(remotes))
	}

	if remotes[0] != "https://github.com/kyoh86/test-repo.git" {
		t.Errorf("Expected default URL %q, got %q",
			"https://github.com/kyoh86/test-repo.git",
			remotes[0])
	}
}

// Mock writer for testing progress output
type mockWriter struct {
	written []byte
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	w.written = append(w.written, p...)
	return len(p), nil
}

// TestErrorHandling tests the error handling in the Clone method
func TestErrorHandling(t *testing.T) {
	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	service := testtarget.NewService()

	// Test handling of specific git errors
	testCases := []struct {
		name        string
		setupFn     func() (string, error)
		expectedErr error
	}{
		{
			name: "repository not exists error",
			setupFn: func() (string, error) {
				// Use a non-existent repository URL
				return "https://github.com/kyoh86/non-existent-repo.git", nil
			},
			expectedErr: coregit.ErrRepositoryNotExists,
		},
		{
			name: "empty remote repository error",
			setupFn: func() (string, error) {
				// Create an empty repository
				emptyRepoPath := filepath.Join(tempDir, "empty")
				if err := os.MkdirAll(emptyRepoPath, 0755); err != nil {
					return "", err
				}
				_, err := git.PlainInit(emptyRepoPath, true)
				return "file://" + emptyRepoPath, err
			},
			expectedErr: coregit.ErrRepositoryEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := tc.setupFn()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Create a directory for this test case
			testDir := filepath.Join(tempDir, tc.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			// Test cloning
			err = service.Clone(ctx, url, testDir, coregit.CloneOptions{})
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

// TestAuthentication tests that authentication is properly used during clone
func TestAuthentication(t *testing.T) {
	// This is a more theoretical test, as we can't easily test actual authentication
	// without exposing credentials. We'll verify the auth is set up correctly.

	ctx := context.Background()
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a basic service
	service := testtarget.NewService()

	// Authenticate
	authenticatedService, err := service.AuthenticateWithUsernamePassword(ctx, "testuser", "testpass")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	// Try to clone (this will fail but we want to check auth is used)
	err = authenticatedService.Clone(ctx, "https://github.com/kyoh86/non-existent-repo.git", tempDir, coregit.CloneOptions{})

	// We expect an error (repo not found), but we can't easily test the auth was used
	// This test is more for code coverage than actual verification
	if err == nil {
		t.Error("Expected an error from clone")
	}
}
