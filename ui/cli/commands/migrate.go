package commands

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(_ context.Context, svc *service.ServiceSet, defaultNameStore *config.DefaultNameStore, tokenStore *config.TokenStore, workspaceStore *config.WorkspaceStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := defaultNameStore.Save(ctx, svc.DefaultNameService, true); err != nil {
				return err
			}
			if err := tokenStore.Save(ctx, svc.TokenService, true); err != nil {
				return err
			}
			if err := workspaceStore.Save(ctx, svc.WorkspaceService, true); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
