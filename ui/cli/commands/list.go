package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/list"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewListCommand(ctx context.Context, svc *service.ServiceSet) *cobra.Command {
	var f config.ListFlags
	var format flags.LocationFormat
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local repositories",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			opts := list.Options{
				Query: f.Query,
				Limit: f.Limit,
			}
			for repo, err := range list.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx, f.Primary, opts) {
				if err != nil {
					log.FromContext(ctx).WithFields(log.Fields{
						"error": err,
					}).Error("failed to list repositories")
					return nil
				}
				str, err := formatter.Format(*repo)
				if err != nil {
					log.FromContext(ctx).WithFields(log.Fields{
						"error":  err,
						"format": format.String(),
						"path":   repo.FullPath(),
					}).Info("failed to format")
				}
				fmt.Println(str)
			}

			return nil
		},
	}
	cmd.Flags().IntVarP(&f.Limit, "limit", "", svc.Flags.List.Limit, "Max number of repositories to list. -1 means unlimited")
	cmd.Flags().StringVarP(&f.Query, "query", "q", "", "Query for selecting repositories")
	cmd.Flags().BoolVarP(&f.Primary, "primary", "", svc.Flags.List.Primary, "List up repositories in just a primary root")
	if err := flags.LocationFormatFlag(cmd, &format, svc.Flags.List.Format); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to init format flag")
	}

	return cmd
}
