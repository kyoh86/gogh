package worktree

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/worktree"
)

// WorktreeWithRepo represents a worktree with its repository information
type WorktreeWithRepo struct {
	Repo     string
	Worktree worktree.Worktree
}

// InitUsecase initializes a worktree use case with services
// The worktreeService parameter should be created by the infrastructure layer
func InitUsecase(worktreeService worktree.Service, workspaceService workspace.WorkspaceService, finderService workspace.FinderService, referenceParser repository.ReferenceParser) *Usecase {
	return NewUsecase(worktreeService, workspaceService, finderService, referenceParser)
}

// Usecase represents the worktree use case
type Usecase struct {
	worktreeService  worktree.Service
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	referenceParser  repository.ReferenceParser
}

// NewUsecase creates a new worktree use case
func NewUsecase(
	worktreeService worktree.Service,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
		worktreeService:  worktreeService,
		workspaceService: workspaceService,
		finderService:    finderService,
		referenceParser:  referenceParser,
	}
}

// ListOptions contains options for the list operation
type ListOptions struct {
	Limit    int
	Patterns []string
	Primary  bool
}

// List lists worktrees
// If repoRef is nil, lists worktrees for all repositories
func (uc *Usecase) List(ctx context.Context, repoRef *string, opts ListOptions) ([]WorktreeWithRepo, error) {
	// If repoRef is specified, list worktrees for that repository only
	if repoRef != nil {
		ref, err := uc.referenceParser.Parse(*repoRef)
		if err != nil {
			return nil, fmt.Errorf("parsing repository reference: %w", err)
		}

		// Find the repository
		var repo *repository.Location
		found := false
		for r, err := range uc.finderService.ListAllRepository(ctx, uc.workspaceService, workspace.ListOptions{}) {
			if err != nil {
				return nil, fmt.Errorf("searching repository: %w", err)
			}
			if r.Ref().String() == ref.String() {
				repo = r
				found = true
				break
			}
		}

		if !found || repo == nil {
			return nil, fmt.Errorf("repository not found: %s", *repoRef)
		}

		worktrees, err := uc.worktreeService.List(ctx, *repo)
		if err != nil {
			return nil, fmt.Errorf("listing worktrees: %w", err)
		}

		result := make([]WorktreeWithRepo, len(worktrees))
		for i, wt := range worktrees {
			result[i] = WorktreeWithRepo{
				Repo:     repo.Ref().String(),
				Worktree: wt,
			}
		}
		return result, nil
	}

	// If repoRef is not specified, list worktrees for all repositories
	var allWorktrees []WorktreeWithRepo

	ws := uc.workspaceService
	if opts.Primary {
		layout := ws.GetLayoutFor(ws.GetPrimaryRoot())
		for repo, err := range uc.finderService.ListRepositoryInRoot(ctx, layout, workspace.ListOptions{
			Limit:    opts.Limit,
			Patterns: opts.Patterns,
		}) {
			if err != nil {
				return nil, fmt.Errorf("searching repository: %w", err)
			}

			worktrees, err := uc.worktreeService.List(ctx, *repo)
			if err != nil {
				return nil, fmt.Errorf("listing worktrees for %s: %w", repo.Ref().String(), err)
			}

			for _, wt := range worktrees {
				allWorktrees = append(allWorktrees, WorktreeWithRepo{
					Repo:     repo.Ref().String(),
					Worktree: wt,
				})
			}
		}
	} else {
		for repo, err := range uc.finderService.ListAllRepository(ctx, ws, workspace.ListOptions{
			Limit:    opts.Limit,
			Patterns: opts.Patterns,
		}) {
			if err != nil {
				return nil, fmt.Errorf("searching repository: %w", err)
			}

			worktrees, err := uc.worktreeService.List(ctx, *repo)
			if err != nil {
				return nil, fmt.Errorf("listing worktrees for %s: %w", repo.Ref().String(), err)
			}

			for _, wt := range worktrees {
				allWorktrees = append(allWorktrees, WorktreeWithRepo{
					Repo:     repo.Ref().String(),
					Worktree: wt,
				})
			}
		}
	}

	return allWorktrees, nil
}

// AddOptions contains options for the add operation
type AddOptions struct {
	CreateBranch bool `yaml:"-" toml:"-"`
}

// Add adds a new worktree
func (uc *Usecase) Add(ctx context.Context, repoRef string, branch string, opts AddOptions) error {
	ref, err := uc.referenceParser.Parse(repoRef)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}

	// Find the repository using ListAllRepository
	var repo *repository.Location
	found := false
	for r, err := range uc.finderService.ListAllRepository(ctx, uc.workspaceService, workspace.ListOptions{}) {
		if err != nil {
			return fmt.Errorf("searching repository: %w", err)
		}
		if r.Ref().String() == ref.String() {
			repo = r
			found = true
			break // Found the repository
		}
	}

	if !found || repo == nil {
		return fmt.Errorf("repository not found: %s", repoRef)
	}

	_, err = uc.worktreeService.Add(ctx, *repo, branch, worktree.AddOptions{
		CreateBranch: opts.CreateBranch,
	})
	if err != nil {
		return fmt.Errorf("adding worktree: %w", err)
	}

	return nil
}

// RemoveOptions contains options for the remove operation
type RemoveOptions struct {
	// Reserved for future use
}

// Remove removes a worktree
func (uc *Usecase) Remove(ctx context.Context, repoRef string, branch string, opts RemoveOptions) error {
	ref, err := uc.referenceParser.Parse(repoRef)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}

	// Find the repository using ListAllRepository
	var repo *repository.Location
	found := false
	for r, err := range uc.finderService.ListAllRepository(ctx, uc.workspaceService, workspace.ListOptions{}) {
		if err != nil {
			return fmt.Errorf("searching repository: %w", err)
		}
		if r.Ref().String() == ref.String() {
			repo = r
			found = true
			break // Found the repository
		}
	}

	if !found || repo == nil {
		return fmt.Errorf("repository not found: %s", repoRef)
	}

	// List worktrees to find the one to remove
	worktrees, err := uc.worktreeService.List(ctx, *repo)
	if err != nil {
		return fmt.Errorf("listing worktrees: %w", err)
	}

	var targetWorktree *worktree.Worktree
	for i := range worktrees {
		if worktrees[i].Branch == branch {
			targetWorktree = &worktrees[i]
			break
		}
	}

	if targetWorktree == nil {
		return fmt.Errorf("worktree not found for branch: %s", branch)
	}

	if err := uc.worktreeService.Remove(ctx, *targetWorktree); err != nil {
		return fmt.Errorf("removing worktree: %w", err)
	}

	return nil
}
