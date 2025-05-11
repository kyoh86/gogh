package clone

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the clone use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	gitService       git.GitService
}

// NewUseCase creates a new clone use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		gitService:       gitService,
	}
}

// Options contains options for the clone operation
type Options struct {
	Alias          *repository.Reference
	TryCloneNotify service.TryCloneNotify
}

// Execute performs the clone operation
func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, opts Options) error {
	// Get repository information from remote
	repo, err := uc.hostingService.GetRepository(ctx, ref)
	if err != nil {
		return err
	}
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService, uc.gitService)
	return repositoryService.TryClone(ctx, repo, ref, opts.Alias, opts.TryCloneNotify)
}
