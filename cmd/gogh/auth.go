package main

import (
	"github.com/kyoh86/gogh/v3/internal/tokenstore"
	"github.com/spf13/cobra"
)

var tokens = tokenstore.TokenManager{}

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
