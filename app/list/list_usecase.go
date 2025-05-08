package list

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase defines the use case for listing repositories
type UseCase struct {
	workspaceService workspace.WorkspaceService
}

// NewUseCase creates a new instance of UseCase
func NewUseCase(
	workspaceService workspace.WorkspaceService,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
	}
}

// Execute retrieves a list of repositories under the specified workspace roots
func (u *UseCase) Execute(ctx context.Context, limit int, primary bool) iter.Seq2[workspace.Repository, error] {
	if primary {
		ws := u.workspaceService
		return ws.GetLayoutFor(ws.GetDefaultRoot()).ListRepository(ctx, limit)
	}
	return u.workspaceService.ListRepository(ctx, limit)
}
