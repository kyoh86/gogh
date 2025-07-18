package fork

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

// Usecase represents the fork use case
type Usecase struct {
	hostingService     hosting.HostingService
	workspaceService   workspace.WorkspaceService
	finderService      workspace.FinderService
	overlayService     overlay.OverlayService
	scriptService      script.ScriptService
	hookService        hook.HookService
	defaultNameService repository.DefaultNameService
	referenceParser    repository.ReferenceParser
	gitService         git.GitService
}

// NewUsecase creates a new fork use case
func NewUsecase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	overlayService overlay.OverlayService,
	scriptService script.ScriptService,
	hookService hook.HookService,
	defaultNameService repository.DefaultNameService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *Usecase {
	return &Usecase{
		hostingService:     hostingService,
		workspaceService:   workspaceService,
		finderService:      finderService,
		overlayService:     overlayService,
		scriptService:      scriptService,
		hookService:        hookService,
		defaultNameService: defaultNameService,
		referenceParser:    referenceParser,
		gitService:         gitService,
	}
}

type HostingOptions = hosting.ForkRepositoryOptions

type TryCloneOptions = try.Options

// Options represents the options for the fork use case
type Options struct {
	TryCloneOptions
	HostingOptions
	Target string
}

func (uc *Usecase) parseRefs(source, target string) (*repository.Reference, *repository.ReferenceWithAlias, error) {
	srcRef, err := uc.referenceParser.Parse(source)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid source: %w", err)
	}
	if target == "" {
		owner, err := uc.defaultNameService.GetDefaultOwnerFor(srcRef.Host())
		if err != nil {
			return nil, nil, fmt.Errorf("getting default owner for %q: %w", srcRef.Host(), err)
		}
		return srcRef, &repository.ReferenceWithAlias{
			Reference: repository.NewReference(srcRef.Host(), owner, srcRef.Name()),
		}, nil
	}
	toRef, err := uc.referenceParser.ParseWithAlias(target)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid target: %w", err)
	}
	if toRef.Reference.Host() != srcRef.Host() {
		return nil, nil, fmt.Errorf("the host of the forked repository must be the same as the original repository")
	}
	if toRef.Reference.Owner() == "" {
		return nil, nil, fmt.Errorf("the owner of the forked repository must be specified")
	}
	return srcRef, toRef, nil
}

// Execute forks a repository and clones it to the local machine
func (uc *Usecase) Execute(ctx context.Context, source string, opts Options) error {
	ref, targetRef, err := uc.parseRefs(source, opts.Target)
	if err != nil {
		return err
	}
	fork, err := uc.hostingService.ForkRepository(ctx, *ref, targetRef.Reference, opts.HostingOptions)
	if err != nil {
		return fmt.Errorf("requesting fork: %w", err)
	}
	tryCloneUsecase := try.NewUsecase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	if err := tryCloneUsecase.Execute(ctx, fork, targetRef.Alias, opts.TryCloneOptions); err != nil {
		return err
	}
	globals := make(map[string]any)
	if fork.Parent != nil {
		globals["parent"] = map[string]any{
			"host":      fork.Parent.Ref.Host(),
			"owner":     fork.Parent.Ref.Owner(),
			"name":      fork.Parent.Ref.Name(),
			"clone_url": fork.Parent.CloneURL,
		}
	}
	if err := invoke.NewUsecase(
		uc.workspaceService,
		uc.finderService,
		uc.hookService,
		uc.overlayService,
		uc.scriptService,
		uc.referenceParser,
	).InvokeForWithGlobals(ctx, invoke.EventPostFork, targetRef.String(), globals); err != nil {
		return fmt.Errorf("invoking hooks after creation: %w", err)
	}
	return nil
}
