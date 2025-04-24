package main

import (
	"github.com/kyoh86/gogh/v3/config"
	"github.com/spf13/cobra"
)

var tokens = config.TokenManager{}

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Manage tokens",
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return saveTokens()
	},
}

func init() {
	configCommand.AddCommand(authCommand)
	facadeCommand.AddCommand(authCommand)
}
