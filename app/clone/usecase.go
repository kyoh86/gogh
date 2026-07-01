package clone

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/apex/log"
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

// Usecase represents the clone use case
type Usecase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	overlayService   overlay.OverlayService
	scriptService    script.ScriptService
	hookService      hook.HookService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
	warnedDirect     map[string]struct{}
	warnedDirectMu   sync.Mutex
}

// NewUsecase creates a new clone use case
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
		warnedDirect:     map[string]struct{}{},
	}
}

type TryCloneOptions = try.Options

// Options contains options for the clone operation
type Options struct {
	TryCloneOptions
}

// Execute performs the clone operation
func (uc *Usecase) Execute(ctx context.Context, refWithAlias string, opts Options) error {
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return err
	}
	// Get repository information from remote
	repo, err := uc.hostingService.GetRepository(ctx, ref.Reference)
	direct := false
	if err != nil {
		metadataErr := err
		repo, err = uc.directRepository(ref.Reference)
		if err != nil {
			return err
		}
		direct = true
		uc.warnDirectFallback(ctx, ref.Reference, repo.CloneURL, metadataErr)
	}
	tryCloneUsecase := try.NewUsecase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	tryOpts := opts.TryCloneOptions
	tryOpts.Direct = direct
	if err := tryCloneUsecase.Execute(ctx, repo, ref.Alias, tryOpts); err != nil {
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
	if err := invoke.NewUsecase(
		uc.workspaceService,
		uc.finderService,
		uc.hookService,
		uc.overlayService,
		uc.scriptService,
		uc.referenceParser,
	).InvokeForWithGlobals(ctx, invoke.EventPostClone, refWithAlias, globals); err != nil {
		return fmt.Errorf("invoking hooks after clone: %w", err)
	}
	return nil
}

func (uc *Usecase) warnDirectFallback(ctx context.Context, ref repository.Reference, cloneURL string, metadataErr error) {
	key := ref.String()
	uc.warnedDirectMu.Lock()
	if _, ok := uc.warnedDirect[key]; ok {
		uc.warnedDirectMu.Unlock()
		return
	}
	uc.warnedDirect[key] = struct{}{}
	uc.warnedDirectMu.Unlock()

	log.FromContext(ctx).Warnf("Repository metadata unavailable for %s; using direct clone: %s", key, cloneURL)
	log.FromContext(ctx).Warn("Skipped metadata setup: upstream remote for fork parent and parent values for post-clone hooks")
	log.FromContext(ctx).Debugf("Repository metadata error for %s: %v", key, metadataErr)
}

func (uc *Usecase) directRepository(ref repository.Reference) (*hosting.Repository, error) {
	u, err := uc.hostingService.GetURLOf(ref)
	if err != nil {
		return nil, fmt.Errorf("building direct clone URL: %w", err)
	}
	cloneURL := strings.TrimSuffix(u.String(), ".git") + ".git"
	return &hosting.Repository{
		Ref:      ref,
		URL:      strings.TrimSuffix(cloneURL, ".git"),
		CloneURL: cloneURL,
	}, nil
}
