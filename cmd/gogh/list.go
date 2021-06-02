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
				for _, project := range projects {
					str, err := f.Format(project)
					if err != nil {
						log.FromContext(ctx).WithFields(log.Fields{
							"error":  err,
							"format": project.FullFilePath(),
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
"full-file-path", "url", "fields" and "fields:[separator]".
`

func init() {
	setup()
	listFlags.Format = defaultFlag.List.Format
	listCommand.Flags().StringVarP(&listFlags.Query, "query", "", "", "Query for selecting projects")
	listCommand.Flags().BoolVarP(&listFlags.Primary, "primary", "", defaultFlag.List.Primary, "List up projects in just a primary root")
	listCommand.Flags().VarP(&listFlags.Format, "format", "f", formatShortUsage)
	facadeCommand.AddCommand(listCommand)
}
