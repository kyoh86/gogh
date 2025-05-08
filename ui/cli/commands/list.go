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

func NewListCommand(conf *config.ConfigStore, defaults *config.FlagStore, workspaceService workspace.WorkspaceService) *cobra.Command {
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
			useCase := list.NewUseCase(workspaceService)
			for repo, err := range useCase.Execute(ctx, 0, f.Primary) {
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
	f.Format = defaults.List.Format
	// TODO: use "query" flag
	cmd.Flags().StringVarP(&f.Query, "query", "q", "", "Query for selecting repositories")
	cmd.Flags().BoolVarP(&f.Primary, "primary", "", defaults.List.Primary, "List up repositories in just a primary root")
	cmd.Flags().VarP(&f.Format, "format", "f", flags.LocalRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteLocalRepoFormat); err != nil {
		panic(err)
	}
	return cmd
}
