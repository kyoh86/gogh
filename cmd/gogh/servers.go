package main

import (
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

var servers gogh.Servers

var serversCommand = &cobra.Command{
	Use:     "servers",
	Short:   "Manage servers",
	Aliases: []string{"server"},
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return saveServers()
	},
}

func init() {
	setup()
	facadeCommand.AddCommand(serversCommand)
}
