package main

import (
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

var tokens gogh.TokenManager

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Manage tokens",
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return saveTokens()
	},
}

func init() {
	facadeCommand.AddCommand(authCommand)
	configCommand.AddCommand(authCommand)
}
