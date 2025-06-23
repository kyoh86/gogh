package cwd

import (
	"context"
	"fmt"
	"os"

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

// Execute retrieves the repository location for the current working directory
func (uc *UseCase) Execute(ctx context.Context) (*repository.Location, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	return uc.finderService.FindByPath(ctx, uc.workspaceService, wd)
}
