package create

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	gitService       git.GitService
}

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

type Options struct {
	Alias          *repository.Reference
	TryCloneNotify service.TryCloneNotify
	hosting.CreateRepositoryOptions
}

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, opts Options) error {
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService, uc.gitService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref, opts.CreateRepositoryOptions)
	if err != nil {
		return err
	}
	return repositoryService.TryClone(ctx, repo, ref, opts.Alias, opts.TryCloneNotify)
}
