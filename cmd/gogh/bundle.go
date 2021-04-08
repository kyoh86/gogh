package main

import (
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var bundleCommand = &cobra.Command{
	Use:   "bundle",
	Short: "Manage bundle",
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return app.SaveServers()
	},
}

func init() {
	facadeCommand.AddCommand(bundleCommand)
}
