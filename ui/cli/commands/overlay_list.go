package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/overlay/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json   bool
		source bool
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered overlays",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.NewUsecase(svc.OverlayService, cmd.OutOrStdout()).Execute(cmd.Context(), f.json, f.source)
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&f.source, "source", "", false, "Output with source code")
	return cmd, nil
}
