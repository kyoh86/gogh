package try_clone

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase provides common operations for repository manipulation
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	overlayService   overlay.OverlayService
	gitService       git.GitService
}

// NewUseCase creates a new instance of RepositoryService.
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	overlayService overlay.OverlayService,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
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
}

// Execute attempts to clone a repository with retry logic.
func (uc *UseCase) Execute(
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
	if err := cloneWithRetry(ctx, gitService, layout, repo.Ref, repo.CloneURL, localPath, opts.Timeout, opts.Notify); err != nil {
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
	timeout time.Duration,
	notify Notify,
) (err error) {
	if notify == nil {
		notify = func(n Status) error { return nil }
	}
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	for {
		toctx, tocancel := context.WithTimeout(ctx, timeout)
		err = gitService.Clone(toctx, cloneURL, localPath, git.CloneOptions{})
		tocancel()
		switch {
		case errors.Is(err, git.ErrRepositoryNotExists), errors.Is(err, context.DeadlineExceeded):
			if err := notify(StatusRetry); err != nil {
				return err
			}
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

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			// next retry
		}
	}
}
