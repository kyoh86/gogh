package main

import (
	"github.com/spf13/cobra"
)

var bundleCommand = &cobra.Command{
	Use:   "bundle",
	Short: "Manage bundle",
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return saveServers()
	},
}

func init() {
	setup()
	facadeCommand.AddCommand(bundleCommand)
}
