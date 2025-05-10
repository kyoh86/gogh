package clone

import (
	"context"

	"github.com/apex/log"
	service "github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"golang.org/x/sync/errgroup"
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

// Options contains options for the clone operation
type Options struct {
	Alias      *repository.Reference
	RetryLimit int
}

// Execute performs the clone operation
func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, opts Options) error {
	// Get repository information from remote
	repo, err := uc.hostingService.GetRepository(ctx, ref)
	if err != nil {
		return err
	}
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService)

	notify := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		empty, err := repositoryService.CloneRepositoryWithRetry(ctx, repo, ref, opts.Alias, opts.RetryLimit, notify)
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
