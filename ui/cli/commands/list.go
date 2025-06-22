package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

// NewListCommand creates a new command to list local repositories.
func NewListCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.ListFlags
	var format flags.LocationFormat
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local repositories",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := config.LocationFormatter(format.String())
			if err != nil {
				return fmt.Errorf("invalid format flag: %w", err)
			}

			ctx := cmd.Context()
			opts := list.Options{
				Primary: f.Primary,
				ListOptions: list.ListOptions{
					Limit:    f.Limit,
					Patterns: f.Patterns,
				},
			}
			for repo, err := range list.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx, opts) {
				if err != nil {
					return fmt.Errorf("listing up repositories: %w", err)
				}
				str, err := formatter.Format(*repo)
				if err != nil {
					log.FromContext(ctx).WithFields(log.Fields{
						"error":  err,
						"format": format.String(),
						"path":   repo.FullPath(),
					}).Info("Failed to format")
				} else {
					fmt.Println(str)
				}
			}

			return nil
		},
	}
	cmd.Flags().IntVarP(&f.Limit, "limit", "", svc.Flags.List.Limit, "Max number of repositories to list. -1 means unlimited")
	cmd.Flags().StringSliceVarP(&f.Patterns, "pattern", "p", nil, "Patterns for selecting repositories")
	cmd.Flags().BoolVarP(&f.Primary, "primary", "", svc.Flags.List.Primary, "List up repositories in just a primary root")
	if err := flags.LocationFormatFlag(cmd, &format, svc.Flags.List.Format); err != nil {
		return nil, fmt.Errorf("initializing format flag: %s", err)
	}
	return cmd, nil
}
