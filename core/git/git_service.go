package git

import "context"

// GitService handles actual Git operations
type GitService interface {
	// Clone performs the actual git clone operation
	Clone(ctx context.Context, remoteURL string, localPath string, options *CloneOptions) error

	// SetRemote configures remote repositories in a git repo
	SetRemotes(ctx context.Context, localPath string, name string, remotes []string) error
	// SetDefaultRemote configures remote repositories in a git repo
	SetDefaultRemotes(ctx context.Context, localPath string, remotes []string) error
}

// CloneOptions contains options for the local clone operation
type CloneOptions struct {
	// Options like recursive, etc.
}
