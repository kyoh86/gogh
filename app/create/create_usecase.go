package create

import (
	"context"
	"fmt"
	"time"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

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

type RepositoryOptions = hosting.CreateRepositoryOptions

// Options contains options for the create operation
type Options struct {
	RequestTimeout time.Duration
	TryCloneNotify service.TryCloneNotify
	RepositoryOptions
}

// Execute performs the create operation
func (uc *UseCase) Execute(ctx context.Context, refWithAlias string, opts Options) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("invalid ref: %w", err)
	}
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService, uc.gitService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref.Reference, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	return repositoryService.TryClone(ctx, repo, ref.Reference, ref.Alias, opts.RequestTimeout, opts.TryCloneNotify)
}
