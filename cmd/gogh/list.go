package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var listFlags struct {
	query  string
	format app.ProjectFormat
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List local projects",
	RunE: func(cmd *cobra.Command, _ []string) error {
		f, err := listFlags.format.Formatter()
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		for _, root := range app.Roots() {
			local := gogh.NewLocalController(root)
			projects, err := local.List(ctx, &gogh.LocalListOption{Query: listFlags.query})
			if err != nil {
				return err
			}
			for _, project := range projects {
				str, err := f(project)
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

const formatShortUsage = `
Print each project in a given format, where <format> can be one of "rel-path", "rel-file-path",
"full-file-path", "url", "fields" and "fields:<separator>".
`

func init() {
	listCommand.Flags().StringVarP(&listFlags.query, "query", "", "", "Query for selecting projects")
	listCommand.Flags().VarP(&listFlags.format, "format", "f", formatShortUsage)
	facadeCommand.AddCommand(listCommand)
}
