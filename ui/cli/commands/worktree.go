package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/service"
	worktreeapp "github.com/kyoh86/gogh/v4/app/worktree"
	"github.com/spf13/cobra"
)

// NewWorktreeCommand creates a new worktree command
func NewWorktreeCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "worktree",
		Short: "Manage git worktrees",
	}

	// worktree list
	listCmd, err := NewWorktreeListCommand(ctx, svc)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(listCmd)

	// worktree add
	addCmd, err := NewWorktreeAddCommand(ctx, svc)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(addCmd)

	// worktree remove
	removeCmd, err := NewWorktreeRemoveCommand(ctx, svc)
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(removeCmd)

	return cmd, nil
}

func initWorktreeUsecase(svc *service.ServiceSet) *worktreeapp.Usecase {
	return worktreeapp.InitUsecase(svc.WorktreeService, svc.WorkspaceService, svc.FinderService, svc.ReferenceParser)
}
