package create

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	overlayStore     overlay.OverlayStore
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	overlayStore overlay.OverlayStore,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		overlayStore:     overlayStore,
		referenceParser:  referenceParser,
		gitService:       gitService,
	}
}

type RepositoryOptions = hosting.CreateRepositoryOptions

type TryCloneOptions = try_clone.Options

// Options contains options for the create operation
type Options struct {
	TryCloneOptions
	RepositoryOptions
}

// Execute performs the create operation
func (uc *UseCase) Execute(ctx context.Context, refWithAlias string, opts Options) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("invalid ref: %w", err)
	}
	repositoryService := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayStore, uc.gitService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref.Reference, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	return repositoryService.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions)
}
