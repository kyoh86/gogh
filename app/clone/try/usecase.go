package try

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase provides common operations for repository manipulation
type Usecase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	overlayService   overlay.OverlayService
	gitService       git.GitService
}

// NewUsecase creates a new instance of RepositoryService.
func NewUsecase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	overlayService overlay.OverlayService,
	gitService git.GitService,
) *Usecase {
	return &Usecase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		overlayService:   overlayService,
		gitService:       gitService,
	}
}

// Status indicates the progress of the TryClone operation.
type Status int

const (
	// StatusEmpty indicates that the repository is empty.
	StatusEmpty Status = iota
	// StatusRetry indicates that the repository does not exist and a retry is needed.
	StatusRetry
)

// Notify is a callback function that is called during the TryClone process.
type Notify func(n Status) error

// RetryLimit is a decorator for TryCloneNotify that limits the number of retries.
func RetryLimit(limit int, notify Notify) Notify {
	if notify == nil {
		notify = func(n Status) error { return nil }
	}
	return func(n Status) error {
		if n == StatusRetry {
			limit--
			if limit < 0 {
				return errors.New("retry limit reached")
			}
		}
		return notify(n)
	}
}

type Options struct {
	// Notify is a callback function to notify the status of the operation.
	Notify Notify
	// Timeout is the maximum wait time for each clone attempt.
	Timeout time.Duration
	// Worktree uses bare + worktree structure (default: true)
	Worktree bool
}

// Execute attempts to clone a repository with retry logic.
func (uc *Usecase) Execute(
	ctx context.Context,
	repo *hosting.Repository,
	alias *repository.Reference,
	opts Options,
) error {
	// Determine local path based on layout
	targetRef := repo.Ref
	if alias != nil {
		targetRef = *alias
	}
	layout := uc.workspaceService.GetPrimaryLayout()
	localPath := layout.PathFor(targetRef)

	// Get the user and token for authentication
	user, token, err := uc.hostingService.GetTokenFor(ctx, repo.Ref.Host(), repo.Ref.Owner())
	if err != nil {
		return err
	}
	gitService, err := uc.gitService.AuthenticateWithUsernamePassword(ctx, user, token.AccessToken)
	if err != nil {
		return err
	}

	// Perform git clone operation
	if err := cloneWithRetry(ctx, gitService, layout, repo.Ref, repo.CloneURL, localPath, opts.Worktree, opts.Timeout, opts.Notify); err != nil {
		return fmt.Errorf("cloning: %w", err)
	}

	// Set up remotes
	if err := gitService.SetDefaultRemotes(ctx, localPath, []string{repo.CloneURL}); err != nil {
		return fmt.Errorf("setting default remote: %w", err)
	}

	// Set up additional remotes if needed
	if repo.Parent != nil {
		if err = gitService.SetRemotes(ctx, localPath, "upstream", []string{repo.Parent.CloneURL}); err != nil {
			return fmt.Errorf("setting upstream remote: %w", err)
		}
	}
	return nil
}

func cloneWithRetry(
	ctx context.Context,
	gitService git.GitService,
	layout workspace.LayoutService,
	ref repository.Reference,
	cloneURL, localPath string,
	worktree bool,
	timeout time.Duration,
	notify Notify,
) error {
	if notify == nil {
		notify = func(n Status) error { return nil }
	}
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// Determine if we should create a bare repository
	cloneCore := cloneWithinTimeout
	if worktree {
		cloneCore = cloneBareWithinTimeout
	}

	for {
		if err := cloneCore(ctx, gitService, cloneURL, localPath, timeout, notify, layout, ref); !errors.Is(err, errContinue) {
			return err
		}
	}
}

var errContinue = errors.New("continue")

func cloneWithinTimeout(
	ctx context.Context,
	gitService git.GitService,
	cloneURL string,
	localPath string,
	timeout time.Duration,
	notify Notify,
	layout workspace.LayoutService,
	ref repository.Reference,
) error {
	toctx, tocancel := context.WithTimeout(ctx, timeout)

	err := gitService.Clone(toctx, cloneURL, localPath, git.CloneOptions{})
	tocancel()
	switch {
	case errors.Is(err, git.ErrRepositoryNotExists), errors.Is(err, context.DeadlineExceeded):
		if err := notify(StatusRetry); err != nil {
			return err
		}
		// continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			// next retry
		}
		return errContinue
	case err == nil:
		return nil
	case errors.Is(err, git.ErrRepositoryEmpty):
		path, err := layout.CreateRepositoryFolder(ref)
		if err != nil {
			return err
		}
		if err := gitService.Init(ctx, cloneURL, path, false, git.InitOptions{}); err != nil {
			return err
		}
		if err := notify(StatusEmpty); err != nil {
			return err
		}
		return nil
	default:
		return err // return immediately for unrecoverable errors
	}
}

func cloneBareWithinTimeout(
	ctx context.Context,
	gitService git.GitService,
	cloneURL string,
	localPath string,
	timeout time.Duration,
	notify Notify,
	layout workspace.LayoutService,
	ref repository.Reference,
) error {
	toctx, tocancel := context.WithTimeout(ctx, timeout)

	err := gitService.Clone(toctx, cloneURL, localPath, git.CloneOptions{
		IsBare: true,
	})
	tocancel()
	switch {
	case errors.Is(err, git.ErrRepositoryNotExists), errors.Is(err, context.DeadlineExceeded):
		if err := notify(StatusRetry); err != nil {
			return err
		}
		// continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			// next retry
		}
		return errContinue
	case err == nil:
		// Create the main worktree

		// Fetch from remote to get remote tracking branches
		if err := gitService.Fetch(ctx, localPath, "origin"); err != nil {
			// Ignore "already up-to-date" errors as they indicate successful clone
			if !errors.Is(err, git.ErrAlreadyUpToDate) {
				return fmt.Errorf("fetching from remote: %w", err)
			}
		}
		// Set the remote HEAD symbolic reference
		if err := gitService.SetRemoteHead(ctx, localPath, "origin"); err != nil {
			return fmt.Errorf("setting remote head: %w", err)
		}
		// Create main branch from remote's default branch (origin/HEAD)
		// Ignore error if branch already exists
		if err := gitService.CreateBranch(ctx, localPath, "main", "origin/HEAD"); err != nil {
			if !errors.Is(err, os.ErrExist) {
				return fmt.Errorf("creating main branch: %w", err)
			}
		}
		// Create .worktree/main worktree
		if err := gitService.AddWorktree(ctx, localPath, "main", ".worktree/main"); err != nil {
			return fmt.Errorf("creating main worktree: %w", err)
		}
		return nil
	case errors.Is(err, git.ErrRepositoryEmpty):
		path, err := layout.CreateRepositoryFolder(ref)
		if err != nil {
			return err
		}
		if err := gitService.Init(ctx, cloneURL, path, true, git.InitOptions{}); err != nil {
			return err
		}
		// Create the main worktree
		// Fetch from remote to get remote tracking branches
		if err := gitService.Fetch(ctx, path, "origin"); err != nil {
			// Ignore "already up-to-date" errors as they indicate successful clone
			if !errors.Is(err, git.ErrAlreadyUpToDate) {
				return fmt.Errorf("fetching from remote: %w", err)
			}
		}
		// Set the remote HEAD symbolic reference
		if err := gitService.SetRemoteHead(ctx, path, "origin"); err != nil {
			return fmt.Errorf("setting remote head: %w", err)
		}
		// Create main branch from remote's default branch (origin/HEAD)
		if err := gitService.CreateBranch(ctx, path, "main", "origin/HEAD"); err != nil {
			return fmt.Errorf("creating main branch: %w", err)
		}
		// Create .worktree/main worktree
		if err := gitService.AddWorktree(ctx, path, "main", ".worktree/main"); err != nil {
			return fmt.Errorf("creating main worktree: %w", err)
		}
		if err := notify(StatusEmpty); err != nil {
			return err
		}
		return nil
	default:
		return err // return immediately for unrecoverable errors
	}
}
