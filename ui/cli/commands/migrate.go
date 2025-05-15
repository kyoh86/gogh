package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate configurations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := svc.DefaultNameStore.Save(ctx, svc.DefaultNameService, true); err != nil {
				return fmt.Errorf("failed to save default names: %w", err)
			}
			if err := svc.TokenStore.Save(ctx, svc.TokenService, true); err != nil {
				return fmt.Errorf("failed to save tokens: %w", err)
			}
			if err := svc.WorkspaceStore.Save(ctx, svc.WorkspaceService, true); err != nil {
				return fmt.Errorf("failed to save workspaces: %w", err)
			}
			return nil
		},
	}
	return cmd
}
