package clone

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/try_clone"
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
	overlayService   workspace.OverlayService
}

// NewUseCase creates a new clone use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
	overlayService workspace.OverlayService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		referenceParser:  referenceParser,
		gitService:       gitService,
		overlayService:   overlayService,
	}
}

type TryCloneOptions = try_clone.Options

// Options contains options for the clone operation
type Options struct {
	TryCloneOptions
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
	repositoryService := try_clone.NewUseCase(
		uc.hostingService,
		uc.workspaceService,
		uc.gitService,
		uc.overlayService,
	)
	if err := repositoryService.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions); err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}
	return nil
}
