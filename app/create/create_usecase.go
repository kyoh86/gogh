package create

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
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
	finderService    workspace.FinderService
	overlayService   overlay.OverlayService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	overlayService overlay.OverlayService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		finderService:    finderService,
		overlayService:   overlayService,
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
	repositoryService := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref.Reference, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	if err := repositoryService.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions); err != nil {
		return fmt.Errorf("cloning: %w", err)
	}
	overlayApplyUseCase := overlay_apply.NewUseCase(
		uc.workspaceService,
		uc.finderService,
		uc.referenceParser,
		uc.overlayService,
	)
	for ov, err := range overlay_find.NewUseCase(
		uc.referenceParser,
		uc.overlayService,
	).Execute(ctx, refWithAlias) {
		if err != nil {
			return fmt.Errorf("finding overlay: %w", err)
		}
		if err := overlayApplyUseCase.Execute(ctx, refWithAlias, ov.RepoPattern, ov.ForInit, ov.RelativePath); err != nil {
			return err
		}
	}
	return nil
}
