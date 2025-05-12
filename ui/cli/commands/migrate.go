package commands

import (
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(svc *ServiceSet, defaultNameStore *config.DefaultNameStore, tokenStore *config.TokenStore, workspaceStore *config.WorkspaceStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := defaultNameStore.Save(ctx, svc.defaultNameService, true); err != nil {
				return err
			}
			if err := tokenStore.Save(ctx, svc.tokenService, true); err != nil {
				return err
			}
			if err := workspaceStore.Save(ctx, svc.workspaceService, true); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
