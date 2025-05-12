package commands

import (
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewAuthCommand(svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
		PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
			// TODO: here...?
			path, err := config.TokensPath()
			if err != nil {
				return err
			}
			store := config.NewTokenStore(path)
			return store.Save(cmd.Context(), svc.tokenService)
		},
	}
}
