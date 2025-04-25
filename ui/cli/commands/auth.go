package commands

import (
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewAuthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return config.SaveTokens()
		},
	}
}
