package commands

import (
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewAuthCommand(tokens auth.TokenService) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
		PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
			path, err := config.TokensPathV0()
			if err != nil {
				return err
			}
			store := config.NewTokenStore(path)
			return store.Save(cmd.Context(), tokens)
		},
	}
}
