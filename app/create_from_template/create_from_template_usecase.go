package create_from_template

import (
	"context"
	"fmt"
	"time"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the use case for creating a repository from a template.
type UseCase struct {
	hostingService   hosting.HostingService
	workspaceService workspace.WorkspaceService
	referenceParser  repository.ReferenceParser
	gitService       git.GitService
}

func NewUseCase(
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
	referenceParser repository.ReferenceParser,
	gitService git.GitService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		workspaceService: workspaceService,
		referenceParser:  referenceParser,
		gitService:       gitService,
	}
}

type RepositoryOptions = hosting.CreateRepositoryFromTemplateOptions

type CreateFromTemplateOptions struct {
	RequestTimeout time.Duration
	TryCloneNotify service.TryCloneNotify
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
	repositoryService := service.NewRepositoryService(uc.hostingService, uc.workspaceService, uc.gitService)
	repo, err := uc.hostingService.CreateRepositoryFromTemplate(ctx, ref.Reference, template, opts.RepositoryOptions)
	if err != nil {
		return fmt.Errorf("creating repository from template: %w", err)
	}
	return repositoryService.TryClone(ctx, repo, ref.Reference, ref.Alias, opts.RequestTimeout, opts.TryCloneNotify)
}
