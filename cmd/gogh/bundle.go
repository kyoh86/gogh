package main

import (
	"github.com/spf13/cobra"
)

var bundleCommand = &cobra.Command{
	Use:   "bundle",
	Short: "Manage bundle",
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return SaveServers()
	},
}

func init() {
	facadeCommand.AddCommand(bundleCommand)
}
