package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	worktreeapp "github.com/kyoh86/gogh/v4/app/worktree"
	"github.com/kyoh86/gogh/v4/core/worktree"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

// NewWorktreeListCommand creates a new worktree list command
func NewWorktreeListCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.ListFlags
	var format flags.WorktreeFormat
	cmd := &cobra.Command{
		Use:   "list [repo-ref]",
		Short: "List worktrees for repositories",
		Long:  flags.WorktreeFormatLongUsage,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktreeUsecase := initWorktreeUsecase(svc)

			var repoRef *string
			if len(args) > 0 {
				repoRef = &args[0]
			}

			opts := worktreeapp.ListOptions{
				Limit:    f.Limit,
				Patterns: f.Patterns,
				Primary:  f.Primary,
			}

			allWorktrees, err := worktreeUsecase.List(cmd.Context(), repoRef, opts)
			if err != nil {
				return fmt.Errorf("listing worktrees: %w", err)
			}

			if len(allWorktrees) == 0 {
				return nil
			}

			formatter, err := worktree.ParseFormat(format.String())
			if err != nil {
				return fmt.Errorf("invalid format flag: %w", err)
			}

			for _, item := range allWorktrees {
				output, err := formatter.Format(item.Worktree, item.Repo)
				if err != nil {
					return fmt.Errorf("formatting worktree: %w", err)
				}
				fmt.Printf("%s\n", output)
			}

			return nil
		},
	}

	if err := flags.WorktreeFormatFlag(cmd, &format, "default"); err != nil {
		return nil, err
	}
	cmd.Flags().IntVar(&f.Limit, "limit", 100, "Max number of repositories to list. -1 means unlimited")
	cmd.Flags().StringSliceVar(&f.Patterns, "pattern", nil, "Patterns for selecting repositories")
	cmd.Flags().BoolVar(&f.Primary, "primary", false, "List up repositories in just a primary root")

	return cmd, nil
}
