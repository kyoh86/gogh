package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewListCommand(conf *config.ConfigStore, defaults *config.FlagStore) *cobra.Command {
	var f config.ListFlags
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local projects",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := f.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			list := conf.GetRoots()
			if f.Primary && len(list) > 1 {
				list = list[0:1]
			}
			for _, root := range list {
				local := gogh.NewLocalController(root)
				projects, err := local.List(ctx, &gogh.LocalListOption{Query: f.Query})
				if err != nil {
					return err
				}
				log.FromContext(ctx).Debugf("found %d projects in %q", len(projects), root)
				for _, project := range projects {
					str, err := formatter.Format(project)
					if err != nil {
						log.FromContext(ctx).WithFields(log.Fields{
							"error":  err,
							"format": f.Format.String(),
							"path":   project.FullFilePath(),
						}).Info("failed to format")
					}
					fmt.Println(str)
				}
			}
			return nil
		},
	}
	f.Format = defaults.List.Format
	cmd.Flags().StringVarP(&f.Query, "query", "q", "", "Query for selecting projects")
	cmd.Flags().BoolVarP(&f.Primary, "primary", "", defaults.List.Primary, "List up projects in just a primary root")
	cmd.Flags().VarP(&f.Format, "format", "f", flags.ProjectFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteProjectFormat); err != nil {
		panic(err)
	}
	return cmd
}
