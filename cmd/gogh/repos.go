package main

import (
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/command"
	"github.com/spf13/cobra"
)

var reposFlags struct {
	query string
}

var reposCommand = &cobra.Command{
	Use:   "repos",
	Short: "List remote repositories",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return command.RemoteList(cmd.Context(), app.Servers(), reposFlags.query)
	},
}

func init() {
	reposCommand.Flags().StringVarP(&listFlags.query, "query", "", "", "Query for selecting projects")
	facadeCommand.AddCommand(reposCommand)
}
