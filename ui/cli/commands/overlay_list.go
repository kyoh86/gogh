package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/overlay_show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewOverlayListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json bool
	}
	width := 60
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List overlays",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				width = w
			}
			var showUseCase interface {
				Execute(ctx context.Context, ov *overlay.Overlay) error
			}
			if f.json {
				showUseCase = overlay_show.NewUseCaseJSON(cmd.OutOrStdout())
			} else {
				showUseCase = overlay_show.NewUseCaseText(cmd.OutOrStdout(), width)
			}
			for overlay, err := range overlay_list.NewUseCase(svc.OverlayService).Execute(ctx) {
				if err != nil {
					return fmt.Errorf("listing overlay: %w", err)
				}
				if err := showUseCase.Execute(ctx, overlay); err != nil {
					return fmt.Errorf("showing overlay %q: %w", overlay.ID(), err)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "j", false, "Output in JSON format")
	return cmd, nil
}
