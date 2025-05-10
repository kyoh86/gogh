package create

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"golang.org/x/sync/errgroup"
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

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, opts Options) error {
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref, opts.CreateRepositoryOptions)
	if err != nil {
		return err
	}
	notify := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		empty, err := repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, opts.Alias, opts.CloneRetryLimit, notify)
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
