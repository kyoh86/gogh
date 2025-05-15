package clone

import (
	"context"
	"time"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the clone use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

// NewUseCase creates a new clone use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		referenceParser:  referenceParser,
		gitService:       gitService,
	}
}

// Options contains options for the clone operation
type Options struct {
	RequestTimeout time.Duration
	TryCloneNotify service.TryCloneNotify
}

// Execute performs the clone operation
func (uc *UseCase) Execute(ctx context.Context, refWithAlias string, opts Options) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return err
	}
	// Get repository information from remote
	repo, err := uc.hostingService.GetRepository(ctx, ref.Reference)
	if err != nil {
		return err
	}
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService, uc.gitService)
	return repositoryService.TryClone(ctx, repo, ref.Reference, ref.Alias, opts.RequestTimeout, opts.TryCloneNotify)
}
