package worktree

import (
	"context"

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
	Add(ctx context.Context, repo repository.Location, branch string) (Worktree, error)

	// Remove a worktree
	Remove(ctx context.Context, worktree Worktree) error

	// Get worktree for current directory
	GetFromPath(ctx context.Context, path string) (*Worktree, error)
}
