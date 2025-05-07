package create_from_template

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

type CreateFromTemplateOptions struct {
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
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	repo, err := uc.hostingService.CreateRepositoryFromTemplate(ctx, ref, template, options.CreateRepositoryFromTemplateOptions)
	if err != nil {
		return err
	}
	return repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, options.Alias, options.CloneRetryLimit)
}
