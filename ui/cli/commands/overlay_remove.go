package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/overlay/remove"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/completion"
	"github.com/spf13/cobra"
)

func NewOverlayRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "remove [flags] <overlay-id>",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove an overlay",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completion.Overlays(cmd.Context(), svc, toComplete)
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			overlayID := args[0]
			if err := remove.NewUsecase(svc.OverlayService).Execute(ctx, overlayID); err != nil {
				return err
			}

			logger.Infof("Removed overlay %s", overlayID)
			return nil
		},
	}
	return cmd, nil
}
