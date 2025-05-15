package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase defines the use case for listing repositories
type UseCase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
}

// NewUseCase creates a new instance of UseCase
func NewUseCase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
) *UseCase {
	return &UseCase{
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
func (u *UseCase) Execute(ctx context.Context, opts Options) iter.Seq2[*repository.Location, error] {
	ws := u.workspaceService
	if opts.Primary {
		layout := ws.GetLayoutFor(ws.GetPrimaryRoot())
		return u.finderService.ListRepositoryInRoot(ctx, layout, opts.ListOptions)
	}
	return u.finderService.ListAllRepository(ctx, ws, opts.ListOptions)
}
