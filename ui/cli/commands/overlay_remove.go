package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/overlay/remove"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "remove [flags] <overlay-id>",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove an overlay",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			overlayID := args[0]
			if err := remove.NewUseCase(svc.OverlayService).Execute(ctx, overlayID); err != nil {
				return err
			}

			logger.Infof("Removed overlay %s", overlayID)
			return nil
		},
	}
	return cmd, nil
}
