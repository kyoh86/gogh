package create

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/kyoh86/gogh/v4/app/hook/invoke"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase represents the create use case
type Usecase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	overlayService   overlay.OverlayService
	scriptService    script.ScriptService
	hookService      hook.HookService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUsecase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	overlayService overlay.OverlayService,
	scriptService script.ScriptService,
	hookService hook.HookService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *Usecase {
	return &Usecase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		finderService:    finderService,
		overlayService:   overlayService,
		scriptService:    scriptService,
		hookService:      hookService,
		referenceParser:  referenceParser,
		gitService:       gitService,
	}
}

type RepositoryOptions = hosting.CreateRepositoryOptions

type TryCloneOptions = try.Options

// Options contains options for the create operation
type Options struct {
	TryCloneOptions
	RepositoryOptions
}

// Execute performs the create operation
func (uc *Usecase) Execute(ctx context.Context, refWithAlias string, opts Options) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("invalid ref: %w", err)
	}
	tryCloneUsecase := try.NewUsecase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	repo, err := uc.hostingService.CreateRepository(ctx, ref.Reference, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	if err := tryCloneUsecase.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions); err != nil {
		return fmt.Errorf("cloning: %w", err)
	}
	if err := invoke.NewUsecase(
		uc.workspaceService,
		uc.finderService,
		uc.hookService,
		uc.overlayService,
		uc.scriptService,
		uc.referenceParser,
	).InvokeFor(ctx, invoke.EventPostCreate, refWithAlias); err != nil {
		return fmt.Errorf("invoking hooks after creation: %w", err)
	}
	return nil
}
