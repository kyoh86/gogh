package main

import (
	"github.com/spf13/cobra"
)

var bundleCommand = &cobra.Command{
	Use:   "bundle",
	Short: "Manage bundle",
}

func init() {
	facadeCommand.AddCommand(bundleCommand)
}
