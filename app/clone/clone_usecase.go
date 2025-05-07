package clone

import (
	"context"

	service "github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the clone use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
}

// NewUseCase creates a new clone use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
	}
}

// CloneOptions contains options for the clone operation
type CloneOptions struct {
	Alias *repository.Reference
}

// Execute performs the clone operation
func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options CloneOptions) error {
	// Get repository information from remote
	repo, err := uc.hostingService.GetRepository(ctx, ref)
	if err != nil {
		return err
	}
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	return repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, options.Alias, 0)
}
