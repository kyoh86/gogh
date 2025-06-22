package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v4/app/overlay_update"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayUpdateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name         string
		relativePath string
		sourcePath   string
	}
	cmd := &cobra.Command{
		Use:   "update [flags] <overlay-id>",
		Short: "Update an existing overlay",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			overlayID := args[0]
			var content io.Reader
			if f.sourcePath != "" {
				c, err := os.Open(f.sourcePath)
				if err != nil {
					return err
				}
				defer c.Close()
				content = c
			}
			if err := overlay_update.NewUseCase(svc.OverlayService).Execute(ctx, overlayID, f.name, f.relativePath, content); err != nil {
				return fmt.Errorf("updating overlay: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the overlay")
	cmd.Flags().StringVar(&f.relativePath, "relative-path", "", "Relative path of the overlay in the repository")
	cmd.Flags().StringVar(&f.sourcePath, "source", "", "Overlay source file path")
	return cmd, nil
}
