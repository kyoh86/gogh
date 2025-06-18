package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		forInit bool
	}
	cmd := &cobra.Command{
		Use:   "add [flags] <name> <target-path> <source-path>",
		Short: "Add an overlay file",
		Args:  cobra.ExactArgs(3),
		Example: `   Add an overlay file to a repository.
   The <name> is the name of the overlay, which is used to identify it.
   The <target-path> is the path where the overlay file will be copied to in the repository.
   The <source-path> is the path to the file you want to add as an overlay.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add vsc-setting /path/to/source/vscode/settings.json .vscode/settings.json

   The overlay file will be copied to the repository when you run ` + "`gogh overlay apply`.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			name := args[0]
			targetPath := args[1]
			sourcePath := args[2]
			if filepath.IsAbs(targetPath) {
				return fmt.Errorf("target path must be relative, got absolute path: %s", targetPath)
			}

			content, err := os.Open(sourcePath)
			if err != nil {
				return err
			}
			defer content.Close()
			id, err := overlay_add.NewUseCase(svc.OverlayService).Execute(ctx, name, targetPath, content)
			if err != nil {
				return err
			}

			logger.Infof("Added overlay file %s -> %s for ID %s", sourcePath, targetPath, id)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Register the overlay for 'gogh create' command")
	return cmd, nil
}
