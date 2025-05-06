package create

import (
	"context"
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5/plumbing/transport"

	git "github.com/go-git/go-git/v5"
	gitcore "github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git"
)

// UseCase represents the create use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
	}
}

type CreateFromTemplateOptions struct {
	Local           bool
	Remote          bool
	Alias           *repository.Reference
	CloneRetryLimit int
	hosting.CreateRepositoryFromTemplateOptions
}

func (uc *UseCase) Execute(
	ctx context.Context,
	ref repository.Reference,
	template repository.Reference,
	options CreateFromTemplateOptions,
) error {
	// TODO: share the processes with CreateRepository
	// Maybe, "Create" is not a valid usecase.
	// It should be splitted to usecases below.
	// - Create remote
	// - Create remote from template
	// - Create and init local
	// - Try to get
	//     - Repeat to get repository while it is not exist.
	//     - When the repository is found that is empty, return that.
	if options.Remote {
		repo, err := uc.hostingService.CreateRepositoryFromTemplate(ctx, ref, template, &options.CreateRepositoryFromTemplateOptions)
		if err != nil {
			return err
		}
		if options.Local {
			// Determine local path based on layout
			targetRef := ref
			if options.Alias != nil {
				targetRef = *options.Alias
			}
			layout := uc.workspaceService.GetDefaultLayout()
			localPath := layout.PathFor(targetRef)

			// Get the user and token for authentication
			user, token, err := uc.hostingService.GetTokenFor(ctx, ref)
			if err != nil {
				return err
			}
			gitService := gitimpl.NewAuthenticatedService(user, token.AccessToken)

			// Perform git clone operation
			for range options.CloneRetryLimit {
				err := gitService.Clone(ctx, repo.CloneURL, localPath, &gitcore.CloneOptions{})
				switch {
				case errors.Is(err, git.ErrRepositoryNotExists) || errors.Is(err, transport.ErrRepositoryNotFound):
					log.FromContext(ctx).Info("waiting the remote repository is ready")
				case errors.Is(err, transport.ErrEmptyRemoteRepository):
					path, err := layout.CreateRepositoryFolder(ref)
					if err != nil {
						return err
					}
					if err := gitimpl.NewService().Init(repo.CloneURL, path, false); err != nil {
						return err
					}
					log.FromContext(ctx).Info("created empty repository")
					return nil
				case err == nil:
					return nil
				default:
					return err
				}
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(1 * time.Second):
				}
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
		}
	} else if options.Local {
		layout := uc.workspaceService.GetDefaultLayout()
		path, err := layout.CreateRepositoryFolder(ref)
		if err != nil {
			return err
		}
		remoteURL, err := uc.hostingService.GetURLOf(ref)
		if err != nil {
			return err
		}
		if err := gitimpl.NewService().Init(remoteURL.String(), path, false); err != nil {
			return err
		}
	}

	return nil
}
