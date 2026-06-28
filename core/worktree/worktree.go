package worktree

import (
	"context"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// Worktree represents a git worktree
type Worktree struct {
	Repository repository.Location
	Branch     string
	Path       string
	Commit     string
}

// Service provides worktree management operations
type Service interface {
	// List all worktrees for a repository
	List(ctx context.Context, repo repository.Location) ([]Worktree, error)

	// Add a new worktree
	Add(ctx context.Context, repo repository.Location, branch string, opts AddOptions) (Worktree, error)

	// Remove a worktree
	Remove(ctx context.Context, worktree Worktree) error

	// Get worktree for current directory
	GetFromPath(ctx context.Context, path string) (*Worktree, error)
}

// AddOptions contains options for adding a worktree
type AddOptions struct {
	CreateBranch bool
}

// GetWorktreePath returns the working directory path for a repository location.
// For worktree structures (bare repository with .wt/<branch>), it returns the worktree path.
// For non-worktree structures, it returns the repository path itself.
//
// Examples:
//   - Worktree structure: "/path/to/github.com/user/repo" -> "/path/to/github.com/user/repo/.wt/main"
//   - Non-worktree structure: "/path/to/github.com/user/repo" -> "/path/to/github.com/user/repo"
func GetWorktreePath(ctx context.Context, repo *repository.Location) (string, error) {
	if repo == nil {
		return "", os.ErrNotExist
	}

	repoPath := repo.FullPath()
	worktreeDir := filepath.Join(repoPath, DirectoryName)

	// Check if .wt directory exists
	if info, err := os.Stat(worktreeDir); err == nil && info.IsDir() {
		// Worktree structure detected
		// Try to find the main branch worktree
		mainWorktreePath := filepath.Join(worktreeDir, "main")
		if info, err := os.Stat(mainWorktreePath); err == nil && info.IsDir() {
			return mainWorktreePath, nil
		}

		// If .wt/main doesn't exist, return the first subdirectory
		entries, err := os.ReadDir(worktreeDir)
		if err == nil && len(entries) > 0 {
			for _, entry := range entries {
				if entry.IsDir() {
					return filepath.Join(worktreeDir, entry.Name()), nil
				}
			}
		}

		// Fallback to repo path if no worktree found
		return repoPath, nil
	}

	// Non-worktree structure
	return repoPath, nil
}
