package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/service"
	worktreeapp "github.com/kyoh86/gogh/v4/app/worktree"
	"github.com/spf13/cobra"
)

// NewWorktreeAddCommand creates a new worktree add command
func NewWorktreeAddCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "add <repo-ref> <branch>",
		Short: "Add a new worktree",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRef := args[0]
			branch := args[1]

			worktreeUsecase := initWorktreeUsecase(svc)
			opts := worktreeapp.AddOptions{}

			logger := log.FromContext(cmd.Context())
			logger.WithField("repo", repoRef).WithField("branch", branch).Info("Adding worktree")

			if err := worktreeUsecase.Add(cmd.Context(), repoRef, branch, opts); err != nil {
				return fmt.Errorf("adding worktree: %w", err)
			}

			logger.Info("Worktree added successfully")
			return nil
		},
	}

	return cmd, nil
}
