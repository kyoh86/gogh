package create

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
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

type Options struct {
	Alias           *repository.Reference
	CloneRetryLimit int
	hosting.CreateRepositoryOptions
}

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options Options) error {
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref, options.CreateRepositoryOptions)
	if err != nil {
		return err
	}
	return repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, options.Alias, options.CloneRetryLimit)
}
