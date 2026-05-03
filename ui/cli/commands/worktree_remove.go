package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/service"
	worktreeapp "github.com/kyoh86/gogh/v4/app/worktree"
	"github.com/spf13/cobra"
)

// NewWorktreeRemoveCommand creates a new worktree remove command
func NewWorktreeRemoveCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "remove <repo-ref> <branch>",
		Short: "Remove a worktree",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRef := args[0]
			branch := args[1]

			worktreeUsecase := initWorktreeUsecase(svc)
			opts := worktreeapp.RemoveOptions{}

			logger := log.FromContext(cmd.Context())
			logger.WithField("repo", repoRef).WithField("branch", branch).Info("Removing worktree")

			if err := worktreeUsecase.Remove(cmd.Context(), repoRef, branch, opts); err != nil {
				return fmt.Errorf("removing worktree: %w", err)
			}

			logger.Info("Worktree removed successfully")
			return nil
		},
	}

	return cmd, nil
}
