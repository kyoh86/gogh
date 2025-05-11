package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/list"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewListCommand(svc *ServiceSet) *cobra.Command {
	var f config.ListFlags
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local repositories",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := f.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			opts := workspace.ListOptions{
				Query: f.Query,
				Limit: f.Limit,
			}
			for repo, err := range list.NewUseCase(svc.workspaceService, svc.finderService).Execute(ctx, f.Primary, opts) {
				if err != nil {
					log.FromContext(ctx).WithFields(log.Fields{
						"error": err,
					}).Error("failed to list repositories")
					return nil
				}
				str, err := formatter.Format(repo)
				if err != nil {
					log.FromContext(ctx).WithFields(log.Fields{
						"error":  err,
						"format": f.Format.String(),
						"path":   repo.FullPath(),
					}).Info("failed to format")
				}
				fmt.Println(str)
			}

			return nil
		},
	}
	f.Format = svc.defaults.List.Format
	cmd.Flags().IntVarP(&f.Limit, "limit", "", DefaultValue(svc.defaults.List.Limit, 100), "Max number of repositories to list. -1 means unlimited")
	cmd.Flags().StringVarP(&f.Query, "query", "q", "", "Query for selecting repositories")
	cmd.Flags().BoolVarP(&f.Primary, "primary", "", svc.defaults.List.Primary, "List up repositories in just a primary root")
	cmd.Flags().VarP(&f.Format, "format", "f", flags.LocalRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteLocalRepoFormat); err != nil {
		panic(err)
	}
	return cmd
}
