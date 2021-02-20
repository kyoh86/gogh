package main

import (
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/command"
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
		return command.LocalList(cmd.Context(), app.Roots(), listFlags.query, f)
	},
}

func init() {
	listCommand.Flags().StringVarP(&listFlags.query, "query", "", "", "Query for selecting projects")
	listCommand.Flags().VarP(&listFlags.format, "format", "f", "Output format")
	facadeCommand.AddCommand(listCommand)
}
