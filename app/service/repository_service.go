package service

import (
	"context"
	"errors"
	"time"

	git "github.com/go-git/go-git/v5"                // TODO: remove this import
	"github.com/go-git/go-git/v5/plumbing/transport" // TODO: remove this import
	gitcore "github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git" // TODO: remove this import; with methodise NewAuthenticatedService infra/git/git_service.go
)

// RepositoryService provides common operations for repository manipulation
type RepositoryService struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
}

func NewRepositoryService(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
) *RepositoryService {
	return &RepositoryService{
		hostingService:   hostingService,
		workspaceService: workspaceService,
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
	gitService := gitimpl.NewAuthenticatedService(user, token.AccessToken)

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

func cloneWithRetry(ctx context.Context, gitService *gitimpl.GitService, layout workspace.LayoutService, ref repository.Reference, cloneURL, localPath string, notify TryCloneNotify) (err error) {
	if notify == nil {
		notify = func(n TryCloneProgress) error { return nil }
	}
	for {
		err = gitService.Clone(ctx, cloneURL, localPath, gitcore.CloneOptions{})
		switch {
		case errors.Is(err, git.ErrRepositoryNotExists) || errors.Is(err, transport.ErrRepositoryNotFound):
			if err := notify(TryCloneProgressRetry); err != nil {
				return err
			}
		case err == nil:
			return nil
		case errors.Is(err, transport.ErrEmptyRemoteRepository):
			path, err := layout.CreateRepositoryFolder(ref)
			if err != nil {
				return err
			}
			if err := gitimpl.NewService().Init(cloneURL, path, false); err != nil {
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
