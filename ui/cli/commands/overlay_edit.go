package commands

import (
	"context"
	"os"

	overlayedit "github.com/kyoh86/gogh/v4/app/overlay/edit"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayEditCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "edit [flags] <overlay-id>",
		Short: "Edit an existing overlay (with $EDITOR)",
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
			// Extract the overlay to a temporary file
			tmpFile, err := os.CreateTemp("", "gogh_overlay_edit_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile.Name())

			if err := overlayedit.NewUsecase(svc.OverlayService).ExtractOverlay(ctx, overlayID, tmpFile); err != nil {
				return err
			}
			tmpFile.Close()

			if err := edit(os.Getenv("EDITOR"), tmpFile.Name()); err != nil {
				return err
			}

			edited, err := os.Open(tmpFile.Name())
			if err != nil {
				return err
			}
			defer edited.Close()
			return overlayedit.NewUsecase(svc.OverlayService).UpdateOverlay(ctx, overlayID, edited)
		},
	}
	return cmd, nil
}
