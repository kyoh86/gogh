package cwd

import (
	"context"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase defines the use case for listing repository locations
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

// Execute retrieves a list of repository locations under the specified workspace roots
func (uc *UseCase) Execute(ctx context.Context, path string) (*repository.Location, error) {
	return uc.finderService.FindByPath(ctx, uc.workspaceService, path)
}
