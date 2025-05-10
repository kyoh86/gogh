package fork

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"golang.org/x/sync/errgroup"
)

// UseCase represents the fork use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
}

// NewUseCase creates a new fork use case
func NewUseCase(hostingService hosting.HostingService) *UseCase {
	return &UseCase{
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
	notify := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		empty, err := repositoryService.CloneRepositoryWithRetry(ctx, fork, target.Reference, target.Alias, opts.CloneRetryLimit, notify)
		if err != nil {
			return err
		}
		if empty {
			log.FromContext(ctx).Info("created empty repository")
		}
		return nil
	})
	eg.Go(func() error {
		for range notify {
			log.FromContext(ctx).Info("waiting the remote repository is ready")
		}
		return nil
	})
	return eg.Wait()
}
