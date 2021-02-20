package main

import (
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/command"
	"github.com/spf13/cobra"
)

var createFlags struct {
	spec   string
	format app.ProjectFormat
}

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return command.Create(cmd.Context(), app.DefaultRoot(), app.Servers(), createFlags.spec, nil, nil)
	},
}

func init() {
	createCommand.Flags().StringVarP(&createFlags.spec, "spec", "", "", "A spec of the repository to create")
	facadeCommand.AddCommand(createCommand)
}
