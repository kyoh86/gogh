package fork

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the fork use case
type UseCase struct {
	gitService       git.GitService
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
}

// NewUseCase creates a new fork use case
func NewUseCase(gitService git.GitService, hostingService hosting.HostingService) *UseCase {
	return &UseCase{
		gitService:     gitService,
		hostingService: hostingService,
	}
}

// Options represents the options for the fork use case
type Options struct {
	CloneRetryLimit int
	hosting.ForkRepositoryOptions
}

// Execute forks a repository and clones it to the local machine
func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, target repository.ReferenceWithAlias, opts Options) error {
	fork, err := uc.hostingService.ForkRepository(ctx, ref, target.Reference, opts.ForkRepositoryOptions)
	if err != nil {
		return err
	}

	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	return repositoryService.CloneRepositoryWithRetry(ctx, fork, target.Reference, target.Alias, opts.CloneRetryLimit)
}
