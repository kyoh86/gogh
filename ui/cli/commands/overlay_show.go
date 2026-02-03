package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/overlay/show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json   bool
		source bool
	}
	cmd := &cobra.Command{
		Use:   "show [flags] <overlay-id>",
		Short: "Show an overlay",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completeOverlays(cmd.Context(), svc, toComplete)
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			overlayID := args[0]
			overlayShowUsecase := show.NewUsecase(svc.OverlayService, cmd.OutOrStdout())
			if err := overlayShowUsecase.Execute(ctx, overlayID, f.json, f.source); err != nil {
				return fmt.Errorf("showing overlay %s: %w", overlayID, err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&f.source, "source", "", false, "Output with source code")
	return cmd, nil
}
