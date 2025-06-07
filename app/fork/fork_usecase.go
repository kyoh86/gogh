package fork

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

// UseCase represents the fork use case
type UseCase struct {
	hostingService     hosting.HostingService
	workspaceService   workspace.WorkspaceService
	overlayService     overlay.OverlayService
	defaultNameService repository.DefaultNameService
	referenceParser    repository.ReferenceParser
	gitService         git.GitService
}

// NewUseCase creates a new fork use case
func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	overlayService overlay.OverlayService,
	defaultNameService repository.DefaultNameService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:     hostingService,
		workspaceService:   workspaceService,
		overlayService:     overlayService,
		defaultNameService: defaultNameService,
		referenceParser:    referenceParser,
		gitService:         gitService,
	}
}

type HostingOptions = hosting.ForkRepositoryOptions

type TryCloneOptions = try_clone.Options

// Options represents the options for the fork use case
type Options struct {
	TryCloneOptions
	HostingOptions
	Target string
}

func (uc *UseCase) parseRefs(source, target string) (*repository.Reference, *repository.ReferenceWithAlias, error) {
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
func (uc *UseCase) Execute(ctx context.Context, source string, opts Options) error {
	ref, targetRef, err := uc.parseRefs(source, opts.Target)
	if err != nil {
		return err
	}
	fork, err := uc.hostingService.ForkRepository(ctx, *ref, targetRef.Reference, opts.HostingOptions)
	if err != nil {
		return fmt.Errorf("requesting fork: %w", err)
	}

	repositoryService := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayService, uc.gitService)
	if err := repositoryService.Execute(ctx, fork, targetRef.Alias, opts.TryCloneOptions); err != nil {
		return fmt.Errorf("cloning forked repository: %w", err)
	}
	return nil
}
