package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/spf13/cobra"
)

func NewOverlayListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json bool
	}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List overlays",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			write := func(ov *overlay.Overlay) {
				fmt.Printf("- ")
				fmt.Printf("Repository pattern: %s\n", ov.RepoPattern)
				if ov.ForInit {
					fmt.Printf("  For Init: Yes\n")
				}
				fmt.Printf("  Overlay path: %s\n", ov.RelativePath)
			}
			if f.json {
				enc := json.NewEncoder(cmd.OutOrStdout())
				write = func(ov *overlay.Overlay) {
					if err := enc.Encode(ov); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error encoding overlay %s: %v\n", ov.ID(), err)
					}
				}
			}
			for overlay, err := range overlay_list.NewUseCase(svc.OverlayService).Execute(ctx) {
				if err != nil {
					return fmt.Errorf("listing overlay: %w", err)
				}
				write(overlay)
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "j", false, "Output in JSON format")
	return cmd, nil
}
