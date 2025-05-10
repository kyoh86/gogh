package service

import (
	"context"
	"errors"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gitcore "github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git"
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

// CloneRepositoryWithRetry attempts to clone a repository with retry logic
func (s *RepositoryService) CloneRepositoryWithRetry(
	ctx context.Context,
	repo *hosting.Repository,
	ref repository.Reference,
	alias *repository.Reference,
	retryLimit int,
	retry chan<- struct{},
) (bool, error) {
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
		return false, err
	}
	gitService := gitimpl.NewAuthenticatedService(user, token.AccessToken)

	// Perform git clone operation
	empty, err := cloneWithRetry(ctx, gitService, layout, ref, repo.CloneURL, localPath, retryLimit, retry)
	if err != nil {
		return empty, err
	}

	// Set up remotes
	if err := gitService.SetDefaultRemotes(ctx, localPath, []string{repo.CloneURL}); err != nil {
		return empty, err
	}

	// Set up additional remotes if needed
	if repo.Parent != nil {
		if err = gitService.SetRemotes(ctx, localPath, "upstream", []string{repo.Parent.CloneURL}); err != nil {
			return empty, err
		}
	}
	return empty, nil
}

func cloneWithRetry(ctx context.Context, gitService *gitimpl.GitService, layout workspace.LayoutService, ref repository.Reference, cloneURL, localPath string, retryLimit int, retry chan<- struct{}) (empty bool, err error) {
	defer close(retry)
	for range retryLimit {
		err = gitService.Clone(ctx, cloneURL, localPath, gitcore.CloneOptions{})
		switch {
		case errors.Is(err, git.ErrRepositoryNotExists) || errors.Is(err, transport.ErrRepositoryNotFound):
			select {
			case retry <- struct{}{}:
			default:
			}
		case err == nil:
			return false, nil
		case errors.Is(err, transport.ErrEmptyRemoteRepository):
			path, err := layout.CreateRepositoryFolder(ref)
			if err != nil {
				return true, err
			}
			if err := gitimpl.NewService().Init(cloneURL, path, false); err != nil {
				return true, err
			}
			return true, nil
		default:
			return false, err // return immediately for unrecoverable errors
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(1 * time.Second):
			// next retry
		}
	}
	return false, err
}
