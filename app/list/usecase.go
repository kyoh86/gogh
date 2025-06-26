package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase defines the use case for listing repositories
type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
}

// NewUsecase creates a new instance of Usecase
func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
) *Usecase {
	return &Usecase{
		workspaceService: workspaceService,
		finderService:    finderService,
	}
}

type ListOptions = workspace.ListOptions

type Options struct {
	Primary bool
	ListOptions
}

// Execute retrieves a list of repositories under the specified workspace roots
func (uc *Usecase) Execute(ctx context.Context, opts Options) iter.Seq2[*repository.Location, error] {
	ws := uc.workspaceService
	if opts.Primary {
		layout := ws.GetLayoutFor(ws.GetPrimaryRoot())
		return uc.finderService.ListRepositoryInRoot(ctx, layout, opts.ListOptions)
	}
	return uc.finderService.ListAllRepository(ctx, ws, opts.ListOptions)
}
