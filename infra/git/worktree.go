package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/worktree"
)

// WorktreeService manages git worktrees
type WorktreeService struct {
	pathBuilder worktree.PathBuilder
}

// NewWorktreeService creates a new WorktreeService
func NewWorktreeService(pathBuilder worktree.PathBuilder) worktree.Service {
	return &WorktreeService{
		pathBuilder: pathBuilder,
	}
}

// List lists all worktrees for a repository
func (s *WorktreeService) List(ctx context.Context, repo repository.Location) ([]worktree.Worktree, error) {
	output, err := s.runGitCommand(ctx, repo.FullPath(), "worktree", "list", "--porcelain")
	if err != nil {
		// If git worktree list fails, it might be a non-bare repository
		// In that case, treat the repository itself as the main worktree
		return []worktree.Worktree{
			{
				Repository: repo,
				Branch:     "", // Will be determined later if needed
				Path:       repo.FullPath(),
				Commit:     "",
			},
		}, nil
	}

	return s.parseListOutput(repo, output)
}

// Add adds a new worktree
func (s *WorktreeService) Add(ctx context.Context, repo repository.Location, branch string, opts worktree.AddOptions) (worktree.Worktree, error) {
	worktreePath := s.pathBuilder.BuildWorktreePath(repo, branch)

	var cmd *exec.Cmd
	// If CreateBranch is requested and branch doesn't exist, use -b flag
	if opts.CreateBranch && !s.branchExists(ctx, repo.FullPath(), branch) {
		// git worktree add -b <new-branch> <path> <start-point>
		cmd = exec.CommandContext(ctx, "git", "-C", repo.FullPath(), "worktree", "add", "-b", branch, worktreePath, "HEAD")
	} else {
		// git worktree add <path> <branch>
		cmd = exec.CommandContext(ctx, "git", "-C", repo.FullPath(), "worktree", "add", worktreePath, branch)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return worktree.Worktree{}, fmt.Errorf("adding worktree: %w\nOutput: %s", err, string(output))
	}

	return worktree.Worktree{
		Repository: repo,
		Branch:     branch,
		Path:       worktreePath,
	}, nil
}

// branchExists checks if a branch exists in the repository
func (s *WorktreeService) branchExists(ctx context.Context, repoPath, branch string) bool {
	cmd := exec.CommandContext(ctx, "git", "-C", repoPath, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	err := cmd.Run()
	return err == nil
}

// Remove removes a worktree
func (s *WorktreeService) Remove(ctx context.Context, wt worktree.Worktree) error {
	// git worktree remove <path>
	cmd := exec.CommandContext(ctx, "git", "-C", wt.Repository.FullPath(), "worktree", "remove", wt.Path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("removing worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetFromPath gets the worktree for a given path
func (s *WorktreeService) GetFromPath(ctx context.Context, path string) (*worktree.Worktree, error) {
	// git worktree list --porcelain to find the worktree containing this path
	cmd := exec.CommandContext(ctx, "git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing worktrees: %w", err)
	}

	// Parse the output to find the worktree
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var currentWorktree *worktree.Worktree
	for scanner.Scan() {
		line := scanner.Text()
		if worktreePath, ok := strings.CutPrefix(line, "worktree "); ok {
			// Check if the given path is within this worktree
			if strings.HasPrefix(path, worktreePath) {
				currentWorktree = &worktree.Worktree{
					Path: worktreePath,
				}
			}
		}
		if currentWorktree != nil && strings.HasPrefix(line, "branch ") {
			branchRef := strings.TrimPrefix(line, "branch ")
			currentWorktree.Branch = strings.TrimPrefix(branchRef, "refs/heads/")
		}
		if line == "" {
			if currentWorktree != nil && strings.HasPrefix(path, currentWorktree.Path) {
				return currentWorktree, nil
			}
			currentWorktree = nil
		}
	}

	return nil, fmt.Errorf("no worktree found for path: %s", path)
}

// runGitCommand executes a git command and returns the output
func (s *WorktreeService) runGitCommand(ctx context.Context, repoPath string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoPath)
	cmd.Args = append(cmd.Args, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running git command: %w", err)
	}
	return string(output), nil
}

// parseListOutput parses "git worktree list --porcelain" output
func (s *WorktreeService) parseListOutput(repo repository.Location, output string) ([]worktree.Worktree, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var worktrees []worktree.Worktree
	var currentWorktree *worktree.Worktree
	isBare := false

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "worktree "):
			worktreePath := strings.TrimPrefix(line, "worktree ")
			currentWorktree = &worktree.Worktree{
				Repository: repo,
				Path:       worktreePath,
			}
			isBare = false
		case strings.HasPrefix(line, "bare"):
			// This is a bare repository, not a worktree
			isBare = true
		case strings.HasPrefix(line, "HEAD "):
			if currentWorktree != nil {
				currentWorktree.Commit = strings.TrimPrefix(line, "HEAD ")
			}
		case strings.HasPrefix(line, "branch "):
			if currentWorktree != nil {
				branchRef := strings.TrimPrefix(line, "branch ")
				currentWorktree.Branch = strings.TrimPrefix(branchRef, "refs/heads/")
			}
		case strings.HasPrefix(line, "detached"):
			// Handle detached HEAD state
		case line == "":
			// Empty line marks the end of a worktree entry
			if currentWorktree != nil && !isBare {
				worktrees = append(worktrees, *currentWorktree)
			}
			currentWorktree = nil
			isBare = false
		}
	}

	// Add the last worktree if there's no trailing empty line
	if currentWorktree != nil && !isBare {
		worktrees = append(worktrees, *currentWorktree)
	}

	return worktrees, nil
}
