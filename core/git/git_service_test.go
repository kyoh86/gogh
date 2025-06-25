package git_test

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/kyoh86/gogh/v4/core/git"
)

func TestGitServiceErrors(t *testing.T) {
	t.Run("ErrRepositoryNotExists", func(t *testing.T) {
		err := git.ErrRepositoryNotExists
		if err.Error() != "repository not exists" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("ErrRepositoryEmpty", func(t *testing.T) {
		err := git.ErrRepositoryEmpty
		if err.Error() != "repository is empty" {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestCloneOptions(t *testing.T) {
	// Test that CloneOptions can be instantiated
	opts := git.CloneOptions{}
	// Currently no fields, but this ensures the struct exists and can be used
	_ = opts
}

func TestInitOptions(t *testing.T) {
	// Test that InitOptions can be instantiated
	opts := git.InitOptions{}
	// Currently no fields, but this ensures the struct exists and can be used
	_ = opts
}

// MockGitService is a mock implementation of GitService for testing
type MockGitService struct {
	AuthenticateFunc      func(ctx context.Context, username, password string) (git.GitService, error)
	CloneFunc             func(ctx context.Context, remoteURL string, localPath string, opts git.CloneOptions) error
	InitFunc              func(ctx context.Context, remoteURL string, localPath string, isBare bool, opts git.InitOptions) error
	SetRemotesFunc        func(ctx context.Context, localPath string, name string, remotes []string) error
	SetDefaultRemotesFunc func(ctx context.Context, localPath string, remotes []string) error
	GetRemotesFunc        func(ctx context.Context, localPath string, name string) ([]string, error)
	GetDefaultRemotesFunc func(ctx context.Context, localPath string) ([]string, error)
	ListExcludedFilesFunc func(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error]
	ListAllFilesFunc      func(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error]
}

func (m *MockGitService) AuthenticateWithUsernamePassword(ctx context.Context, username, password string) (git.GitService, error) {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(ctx, username, password)
	}
	return m, nil
}

func (m *MockGitService) Clone(ctx context.Context, remoteURL string, localPath string, opts git.CloneOptions) error {
	if m.CloneFunc != nil {
		return m.CloneFunc(ctx, remoteURL, localPath, opts)
	}
	return nil
}

func (m *MockGitService) Init(ctx context.Context, remoteURL string, localPath string, isBare bool, opts git.InitOptions) error {
	if m.InitFunc != nil {
		return m.InitFunc(ctx, remoteURL, localPath, isBare, opts)
	}
	return nil
}

func (m *MockGitService) SetRemotes(ctx context.Context, localPath string, name string, remotes []string) error {
	if m.SetRemotesFunc != nil {
		return m.SetRemotesFunc(ctx, localPath, name, remotes)
	}
	return nil
}

func (m *MockGitService) SetDefaultRemotes(ctx context.Context, localPath string, remotes []string) error {
	if m.SetDefaultRemotesFunc != nil {
		return m.SetDefaultRemotesFunc(ctx, localPath, remotes)
	}
	return nil
}

func (m *MockGitService) GetRemotes(ctx context.Context, localPath string, name string) ([]string, error) {
	if m.GetRemotesFunc != nil {
		return m.GetRemotesFunc(ctx, localPath, name)
	}
	return []string{}, nil
}

func (m *MockGitService) GetDefaultRemotes(ctx context.Context, localPath string) ([]string, error) {
	if m.GetDefaultRemotesFunc != nil {
		return m.GetDefaultRemotesFunc(ctx, localPath)
	}
	return []string{}, nil
}

func (m *MockGitService) ListExcludedFiles(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error] {
	if m.ListExcludedFilesFunc != nil {
		return m.ListExcludedFilesFunc(ctx, localPath, filePatterns)
	}
	return func(yield func(string, error) bool) {}
}

func (m *MockGitService) ListAllFiles(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error] {
	if m.ListAllFilesFunc != nil {
		return m.ListAllFilesFunc(ctx, localPath, filePatterns)
	}
	return func(yield func(string, error) bool) {}
}

// TestGitServiceInterface verifies that MockGitService implements GitService
func TestGitServiceInterface(t *testing.T) {
	// This test will fail to compile if MockGitService doesn't implement GitService
	var _ git.GitService = (*MockGitService)(nil)
}

// TestGitServiceMethods tests that each method can be called and returns expected results
func TestGitServiceMethods(t *testing.T) {
	ctx := context.Background()

	t.Run("AuthenticateWithUsernamePassword", func(t *testing.T) {
		expectedError := errors.New("auth failed")
		mock := &MockGitService{
			AuthenticateFunc: func(ctx context.Context, username, password string) (git.GitService, error) {
				if username != "testuser" || password != "testpass" {
					return nil, expectedError
				}
				return &MockGitService{}, nil
			},
		}

		// Test successful authentication
		result, err := mock.AuthenticateWithUsernamePassword(ctx, "testuser", "testpass")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}

		// Test failed authentication
		_, err = mock.AuthenticateWithUsernamePassword(ctx, "wronguser", "wrongpass")
		if err != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("Clone", func(t *testing.T) {
		expectedError := git.ErrRepositoryNotExists
		mock := &MockGitService{
			CloneFunc: func(ctx context.Context, remoteURL string, localPath string, opts git.CloneOptions) error {
				if remoteURL == "invalid" {
					return expectedError
				}
				return nil
			},
		}

		// Test successful clone
		err := mock.Clone(ctx, "https://github.com/user/repo.git", "/tmp/repo", git.CloneOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Test failed clone
		err = mock.Clone(ctx, "invalid", "/tmp/repo", git.CloneOptions{})
		if err != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("Init", func(t *testing.T) {
		expectedError := errors.New("init failed")
		mock := &MockGitService{
			InitFunc: func(ctx context.Context, remoteURL string, localPath string, isBare bool, opts git.InitOptions) error {
				if localPath == "/invalid/path" {
					return expectedError
				}
				return nil
			},
		}

		// Test successful init
		err := mock.Init(ctx, "https://github.com/user/repo.git", "/tmp/repo", false, git.InitOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Test bare repository init
		err = mock.Init(ctx, "https://github.com/user/repo.git", "/tmp/bare", true, git.InitOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Test failed init
		err = mock.Init(ctx, "https://github.com/user/repo.git", "/invalid/path", false, git.InitOptions{})
		if err != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, err)
		}
	})

	t.Run("SetRemotes", func(t *testing.T) {
		callCount := 0
		mock := &MockGitService{
			SetRemotesFunc: func(ctx context.Context, localPath string, name string, remotes []string) error {
				callCount++
				if name == "upstream" && len(remotes) == 2 {
					return nil
				}
				return errors.New("invalid parameters")
			},
		}

		// Test setting remotes
		err := mock.SetRemotes(ctx, "/tmp/repo", "upstream", []string{"url1", "url2"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected SetRemotesFunc to be called once, got %d", callCount)
		}
	})

	t.Run("SetDefaultRemotes", func(t *testing.T) {
		callCount := 0
		mock := &MockGitService{
			SetDefaultRemotesFunc: func(ctx context.Context, localPath string, remotes []string) error {
				callCount++
				return nil
			},
		}

		// Test setting default remotes
		err := mock.SetDefaultRemotes(ctx, "/tmp/repo", []string{"https://github.com/user/repo.git"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected SetDefaultRemotesFunc to be called once, got %d", callCount)
		}
	})

	t.Run("GetRemotes", func(t *testing.T) {
		expectedRemotes := []string{"url1", "url2"}
		mock := &MockGitService{
			GetRemotesFunc: func(ctx context.Context, localPath string, name string) ([]string, error) {
				if name == "upstream" {
					return expectedRemotes, nil
				}
				return nil, errors.New("remote not found")
			},
		}

		// Test getting existing remote
		remotes, err := mock.GetRemotes(ctx, "/tmp/repo", "upstream")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(remotes) != len(expectedRemotes) {
			t.Errorf("expected %d remotes, got %d", len(expectedRemotes), len(remotes))
		}

		// Test getting non-existent remote
		_, err = mock.GetRemotes(ctx, "/tmp/repo", "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent remote")
		}
	})

	t.Run("GetDefaultRemotes", func(t *testing.T) {
		expectedRemotes := []string{"https://github.com/user/repo.git"}
		mock := &MockGitService{
			GetDefaultRemotesFunc: func(ctx context.Context, localPath string) ([]string, error) {
				return expectedRemotes, nil
			},
		}

		// Test getting default remotes
		remotes, err := mock.GetDefaultRemotes(ctx, "/tmp/repo")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(remotes) != len(expectedRemotes) {
			t.Errorf("expected %d remotes, got %d", len(expectedRemotes), len(remotes))
		}
	})

	t.Run("ListExcludedFiles", func(t *testing.T) {
		expectedFiles := []string{"/tmp/repo/.gitignore", "/tmp/repo/build/"}
		mock := &MockGitService{
			ListExcludedFilesFunc: func(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error] {
				return func(yield func(string, error) bool) {
					for _, file := range expectedFiles {
						if !yield(file, nil) {
							return
						}
					}
				}
			},
		}

		// Test listing excluded files
		var files []string
		for file, err := range mock.ListExcludedFiles(ctx, "/tmp/repo", []string{"*.log"}) {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				break
			}
			files = append(files, file)
		}
		if len(files) != len(expectedFiles) {
			t.Errorf("expected %d files, got %d", len(expectedFiles), len(files))
		}
	})

	t.Run("ListAllFiles", func(t *testing.T) {
		expectedFiles := []string{"/tmp/repo/main.go", "/tmp/repo/README.md", "/tmp/repo/.git/HEAD"}
		mock := &MockGitService{
			ListAllFilesFunc: func(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error] {
				return func(yield func(string, error) bool) {
					for _, file := range expectedFiles {
						if !yield(file, nil) {
							return
						}
					}
				}
			},
		}

		// Test listing all files
		var files []string
		for file, err := range mock.ListAllFiles(ctx, "/tmp/repo", nil) {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				break
			}
			files = append(files, file)
		}
		if len(files) != len(expectedFiles) {
			t.Errorf("expected %d files, got %d", len(expectedFiles), len(files))
		}
	})

	t.Run("ListFilesWithError", func(t *testing.T) {
		expectedError := errors.New("permission denied")
		mock := &MockGitService{
			ListExcludedFilesFunc: func(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error] {
				return func(yield func(string, error) bool) {
					yield("", expectedError)
				}
			},
		}

		// Test error handling in iterator
		var gotError error
		for _, err := range mock.ListExcludedFiles(ctx, "/tmp/repo", nil) {
			if err != nil {
				gotError = err
				break
			}
		}
		if gotError != expectedError {
			t.Errorf("expected error %v, got %v", expectedError, gotError)
		}
	})
}
