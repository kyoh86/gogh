package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

type listFlagsStruct struct {
	Query   string        `yaml:"-"`
	Format  ProjectFormat `yaml:"format,omitempty"`
	Primary bool          `yaml:"primary,omitempty"`
}

var (
	listFlags   listFlagsStruct
	listCommand = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local projects",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := listFlags.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			list := roots()
			if listFlags.Primary && len(list) > 1 {
				list = list[0:1]
			}
			for _, root := range list {
				local := gogh.NewLocalController(root)
				projects, err := local.List(ctx, &gogh.LocalListOption{Query: listFlags.Query})
				if err != nil {
					return err
				}
				log.FromContext(ctx).Debugf("found %d projects in %q", len(projects), root)
				for _, project := range projects {
					str, err := f.Format(project)
					if err != nil {
						log.FromContext(ctx).WithFields(log.Fields{
							"error":  err,
							"format": listFlags.Format.String(),
							"path":   project.FullFilePath(),
						}).Info("failed to format")
					}
					fmt.Println(str)
				}
			}
			return nil
		},
	}
)

const formatShortUsage = `
Print each project in a given format, where [format] can be one of "rel-path", "rel-file-path",
"full-file-path", "json", "url", "fields" or "fields:[separator]".
`

func completeFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"rel-path", "rel-file-path", "full-file-path", "json", "url", "fields", "fields:"}, cobra.ShellCompDirectiveDefault
}

func init() {
	listFlags.Format = defaultFlag.List.Format
	listCommand.Flags().
		StringVarP(&listFlags.Query, "query", "q", "", "Query for selecting projects")
	listCommand.Flags().
		BoolVarP(&listFlags.Primary, "primary", "", defaultFlag.List.Primary, "List up projects in just a primary root")
	listCommand.Flags().VarP(&listFlags.Format, "format", "f", formatShortUsage)
	if err := listCommand.RegisterFlagCompletionFunc("format", completeFormat); err != nil {
		panic(err)
	}
	facadeCommand.AddCommand(listCommand)
}
