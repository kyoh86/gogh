package clone

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

// UseCase represents the clone use case
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

// NewUseCase creates a new clone use case
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
	tryCloneUseCase := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	if err := tryCloneUseCase.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions); err != nil {
		return err
	}
	globals := make(map[string]any)
	if repo.Parent != nil {
		globals["parent"] = map[string]any{
			"host":      repo.Parent.Ref.Host(),
			"owner":     repo.Parent.Ref.Owner(),
			"name":      repo.Parent.Ref.Name(),
			"clone_url": repo.Parent.CloneURL,
		}
	}
	if err := hook_invoke.NewUseCase(
		uc.workspaceService,
		uc.finderService,
		uc.hookService,
		uc.overlayService,
		uc.scriptService,
		uc.referenceParser,
	).InvokeForWithGlobals(ctx, hook_invoke.EventPostClone, refWithAlias, globals); err != nil {
		return fmt.Errorf("invoking hooks after clone: %w", err)
	}
	return nil
}
