package service

import (
	"context"
	"errors"
	"time"

	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// RepositoryService provides common operations for repository manipulation
type RepositoryService struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	gitService       git.GitService
}

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

type TryCloneProgress int

const (
	TryCloneProgressEmpty TryCloneProgress = iota
	TryCloneProgressRetry
)

type TryCloneNotify func(n TryCloneProgress) error

func RetryLimit(limit int, notify TryCloneNotify) TryCloneNotify {
	if notify == nil {
		notify = func(n TryCloneProgress) error { return nil }
	}
	return func(n TryCloneProgress) error {
		if n == TryCloneProgressRetry {
			limit--
			if limit <= 0 {
				return errors.New("retry limit reached")
			}
		}
		return notify(n)
	}
}

func (s *RepositoryService) TryClone(
	ctx context.Context,
	repo *hosting.Repository,
	ref repository.Reference,
	alias *repository.Reference,
	notify TryCloneNotify,
) error {
	// TryClone attempts to clone a repository with retry logic.
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
	if err := cloneWithRetry(ctx, gitService, layout, ref, repo.CloneURL, localPath, notify); err != nil {
		return err
	}

	// Set up remotes
	if err := gitService.SetDefaultRemotes(ctx, localPath, []string{repo.CloneURL}); err != nil {
		return err
	}

	// Set up additional remotes if needed
	if repo.Parent != nil {
		if err = gitService.SetRemotes(ctx, localPath, "upstream", []string{repo.Parent.CloneURL}); err != nil {
			return err
		}
	}
	return nil
}

func cloneWithRetry(ctx context.Context, gitService git.GitService, layout workspace.LayoutService, ref repository.Reference, cloneURL, localPath string, notify TryCloneNotify) (err error) {
	if notify == nil {
		notify = func(n TryCloneProgress) error { return nil }
	}
	for {
		err = gitService.Clone(ctx, cloneURL, localPath, git.CloneOptions{})
		switch {
		case errors.Is(err, git.ErrRepositoryNotExists):
			if err := notify(TryCloneProgressRetry); err != nil {
				return err
			}
		case err == nil:
			return nil
		case errors.Is(err, git.ErrRepositoryEmpty):
			path, err := layout.CreateRepositoryFolder(ref)
			if err != nil {
				return err
			}
			if err := gitService.Init(cloneURL, path, false); err != nil {
				return err
			}
			if err := notify(TryCloneProgressEmpty); err != nil {
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
