package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
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

// Execute retrieves a list of repositories under the specified workspace roots
func (u *UseCase) Execute(ctx context.Context, primary bool, opts workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	ws := u.workspaceService
	if primary {
		layout := ws.GetLayoutFor(ws.GetPrimaryRoot())
		return u.finderService.ListRepositoryInRoot(ctx, layout, opts)
	}
	return u.finderService.ListAllRepository(ctx, ws, opts)
}
