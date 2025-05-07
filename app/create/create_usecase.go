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

type CreateOptions struct {
	Local           bool
	Remote          bool
	Alias           *repository.Reference
	CloneRetryLimit int
	hosting.CreateRepositoryOptions
}

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options CreateOptions) error {
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	if options.Remote {
		repo, err := uc.hostingService.CreateRepository(ctx, ref, options.CreateRepositoryOptions)
		if err != nil {
			return err
		}
		if options.Local {
			return repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, options.Alias, options.CloneRetryLimit)
		}
	} else if options.Local {
		repositoryService.CreateLocalRepository(ctx, ref)
	}

	return nil
}
