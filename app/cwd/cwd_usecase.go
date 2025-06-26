package cwd

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase defines the use case for listing repository locations
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

// Execute retrieves the repository location for the current working directory
func (uc *Usecase) Execute(ctx context.Context) (*repository.Location, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}
	return uc.finderService.FindByPath(ctx, uc.workspaceService, wd)
}
