package create_from_template

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the use case for creating a repository from a template.
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	overlayStore     workspace.OverlayStore
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	overlayStore workspace.OverlayStore,
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
	repositoryService := try_clone.NewUseCase(uc.hostingService, uc.workspaceService, uc.overlayStore, uc.gitService)
	repo, err := uc.hostingService.CreateRepositoryFromTemplate(ctx, ref.Reference, template, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating repository from template: %w", err)
	}
	return repositoryService.Execute(ctx, repo, ref.Alias, opts.TryCloneOptions)
}
