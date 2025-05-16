package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// RepositoryService provides common operations for repository manipulation
type RepositoryService struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	gitService       git.GitService
}

// NewRepositoryService creates a new instance of RepositoryService.
func NewRepositoryService(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	gitService git.GitService,
) *RepositoryService {
	return &RepositoryService{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		gitService:       gitService,
	}
}

// TryCloneStatus indicates the progress of the TryClone operation.
type TryCloneStatus int

const (
	// TryCloneStatusEmpty indicates that the repository is empty.
	TryCloneStatusEmpty TryCloneStatus = iota
	// TryCloneStatusRetry indicates that the repository does not exist and a retry is needed.
	TryCloneStatusRetry
)

// TryCloneNotify is a callback function that is called during the TryClone process.
type TryCloneNotify func(n TryCloneStatus) error

// RetryLimit is a decorator for TryCloneNotify that limits the number of retries.
func RetryLimit(limit int, notify TryCloneNotify) TryCloneNotify {
	if notify == nil {
		notify = func(n TryCloneStatus) error { return nil }
	}
	return func(n TryCloneStatus) error {
		if n == TryCloneStatusRetry {
			limit--
			if limit <= 0 {
				return errors.New("retry limit reached")
			}
		}
		return notify(n)
	}
}

// TryClone attempts to clone a repository with retry logic.
func (s *RepositoryService) TryClone(
	ctx context.Context,
	repo *hosting.Repository,
	ref repository.Reference,
	alias *repository.Reference,
	requestTimeout time.Duration,
	notify TryCloneNotify,
) error {
	// Determine local path based on layout
	targetRef := ref
	if alias != nil {
		targetRef = *alias
	}
	layout := s.workspaceService.GetPrimaryLayout()
	localPath := layout.PathFor(targetRef)

	// Get the user and token for authentication
	user, token, err := s.hostingService.GetTokenFor(ctx, ref)
	if err != nil {
		return err
	}
	gitService, err := s.gitService.AuthenticateWithUsernamePassword(ctx, user, token.AccessToken)
	if err != nil {
		return err
	}

	// Perform git clone operation
	if err := cloneWithRetry(ctx, gitService, layout, ref, repo.CloneURL, localPath, requestTimeout, notify); err != nil {
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
	notify TryCloneNotify,
) (err error) {
	if notify == nil {
		notify = func(n TryCloneStatus) error { return nil }
	}
	for {
		if timeout == 0 {
			timeout = 5 * time.Second
		}
		toctx, tocancel := context.WithTimeout(ctx, timeout)
		err = gitService.Clone(toctx, cloneURL, localPath, git.CloneOptions{})
		tocancel()
		switch {
		case errors.Is(err, git.ErrRepositoryNotExists), errors.Is(err, context.DeadlineExceeded):
			if err := notify(TryCloneStatusRetry); err != nil {
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
			if err := notify(TryCloneStatusEmpty); err != nil {
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
