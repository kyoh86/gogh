package create_from_template

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_invoke"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the use case for creating a repository from a template.
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	overlayService   overlay.OverlayService
	scriptService    script.ScriptService
	hookService      hook.HookService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	overlayService overlay.OverlayService,
	scriptService script.ScriptService,
	hookService hook.HookService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
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

type RepositoryOptions = hosting.CreateRepositoryFromTemplateOptions

type TryCloneOptions = try_clone.Options

type CreateFromTemplateOptions struct {
	TryCloneOptions
	RepositoryOptions
}

func (uc *UseCase) Execute(
	ctx context.Context,
	refWithAlias string,
	template repository.Reference,
	opts CreateFromTemplateOptions,
) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("invalid reference: %w", err)
	}
	tryCloneUseCase := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	repo, err := uc.hostingService.CreateRepositoryFromTemplate(ctx, ref.Reference, template, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating repository from template: %w", err)
	}
	if err := tryCloneUseCase.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions); err != nil {
		return fmt.Errorf("cloning: %w", err)
	}
	if err := hook_invoke.NewUseCase(
		uc.workspaceService,
		uc.finderService,
		uc.hookService,
		uc.overlayService,
		uc.scriptService,
		uc.referenceParser,
	).InvokeFor(ctx, hook_invoke.EventPostCreate, refWithAlias); err != nil {
		return fmt.Errorf("invoking hooks after creation: %w", err)
	}
	return nil
}
