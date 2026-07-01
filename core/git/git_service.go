package git

import (
	"context"
	"errors"
	"iter"
)

var (
	// ErrRepositoryNotExists is returned when the repository does not exist
	ErrRepositoryNotExists = errors.New("repository not exists")

	// ErrRepositoryEmpty is returned when the repository is empty
	ErrRepositoryEmpty = errors.New("repository is empty")

	// ErrAlreadyUpToDate is returned when the remote is up-to-date
	ErrAlreadyUpToDate = errors.New("already up-to-date")
)

// GitService handles actual Git operations
type GitService interface {
	// AuthenticateWithUsernamePassword authenticates with a username and password
	AuthenticateWithUsernamePassword(ctx context.Context, username, password string) (GitService, error)

	// Clone performs the actual git clone operation
	Clone(ctx context.Context, remoteURL string, localPath string, opts CloneOptions) error

	// Init initializes a new git repository at the specified local path
	Init(ctx context.Context, remoteURL string, localPath string, isBare bool, opts InitOptions) error

	// IsBare checks if a repository at the given path is a bare repository
	IsBare(ctx context.Context, localPath string) (bool, error)

	// AddWorktree creates a new worktree for an existing repository
	AddWorktree(ctx context.Context, repoPath string, branch string, path string) error

	// Fetch fetches updates from a remote repository
	Fetch(ctx context.Context, repoPath string, remote string) error

	// EnsureRemoteFetchRefspec configures a remote to create remote-tracking refs when fetched
	EnsureRemoteFetchRefspec(ctx context.Context, repoPath string, remote string) error

	// SetRemoteHead sets the HEAD symbolic reference for a remote
	SetRemoteHead(ctx context.Context, repoPath string, remote string) error

	// CreateBranch creates a new branch from a starting point
	CreateBranch(ctx context.Context, repoPath string, branchName string, startPoint string) error

	// SetRemote configures remote repositories in a git repo
	SetRemotes(ctx context.Context, localPath string, name string, remotes []string) error
	// SetDefaultRemote configures the default remote repositories (for usually 'origin') in a git repo
	SetDefaultRemotes(ctx context.Context, localPath string, remotes []string) error

	// GetRemotes retrieves remote repositories from a git repo
	GetRemotes(
		ctx context.Context,
		localPath string,
		name string,
	) ([]string, error)

	// GetDefaultRemotes retrieves the default remote repositories (for usually 'origin') from a git repo
	GetDefaultRemotes(
		ctx context.Context,
		localPath string,
	) ([]string, error)

	// ListExcludedFiles returns a list of untracked files in the repository
	ListExcludedFiles(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error]

	// ListAllFiles returns a list of untracked files in the repository
	ListAllFiles(ctx context.Context, localPath string, filePatterns []string) iter.Seq2[string, error]
}

// CloneOptions contains options for the local clone operation
type CloneOptions struct {
	// IsBare specifies whether to create a bare repository
	IsBare bool
	// UseSystemGit clones with the git executable instead of the default implementation.
	UseSystemGit bool
	// Reserved for future use
}

// InitOptions contains options for the local clone operation
type InitOptions struct {
	// Reserved for future use
}
