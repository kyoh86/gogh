package worktree

import (
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// PathBuilder generates worktree paths
type PathBuilder interface {
	BuildWorktreePath(repo repository.Location, branch string) string
}

type pathBuilder struct{}

// NewPathBuilder creates a new PathBuilder
func NewPathBuilder() PathBuilder {
	return &pathBuilder{}
}

// BuildWorktreePath generates the path for a worktree
// It uses .worktree/ subdirectory and preserves branch name slashes as subdirectories
// Example: repo.Location.FullPath() = "/home/user/Projects/github.com/user/repo"
//
//	branch = "feature/auth"
//	returns "/home/user/Projects/github.com/user/repo/.worktree/feature/auth"
func (b *pathBuilder) BuildWorktreePath(repo repository.Location, branch string) string {
	return filepath.Join(repo.FullPath(), ".worktree", branch)
}
