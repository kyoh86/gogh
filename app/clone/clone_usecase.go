package clone

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git"
)

// UseCase represents the clone use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	layout           workspace.Layout
}

// NewUseCase creates a new clone use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	layout workspace.Layout,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		layout:           layout,
	}
}

// CloneOptions contains options for the clone operation
type CloneOptions struct {
	Alias *repository.Reference
}

// Execute performs the clone operation
func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options *CloneOptions) error {
	// Get repository information from remote
	remoteRepo, err := uc.hostingService.GetRepository(ctx, ref)
	if err != nil {
		return err
	}

	// Determine local path based on layout
	root := uc.workspaceService.GetDefaultRoot()
	targetRef := ref
	if options != nil && options.Alias != nil {
		targetRef = *options.Alias
	}
	localPath := uc.layout.PathFor(root, &targetRef)

	// Get the user and token for authentication
	user, token, err := uc.hostingService.GetTokenFor(ctx, ref)
	if err != nil {
		return err
	}
	gitService := gitimpl.NewAuthenticatedService(user, token.AccessToken)

	// Perform git clone operation
	if err := gitService.Clone(ctx, remoteRepo.CloneURL, localPath, &git.CloneOptions{}); err != nil {
		return err
	}

	// Set up remotes
	if err := gitService.SetDefaultRemotes(ctx, localPath, []string{remoteRepo.CloneURL}); err != nil {
		return err
	}

	// Set up additional remotes if needed
	if remoteRepo.Parent != nil {
		if err = gitService.SetRemotes(ctx, localPath, "upstream", []string{remoteRepo.Parent.CloneURL}); err != nil {
			return err
		}
	}

	return nil
}
