package cwd

import (
	"context"

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
func (u *UseCase) Execute(ctx context.Context, path string) (workspace.RepoInfo, error) {
	return u.finderService.FindByPath(ctx, u.workspaceService, path)
}
