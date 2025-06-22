package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v4/app/script_update"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptUpdateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name       string
		sourcePath string
	}
	cmd := &cobra.Command{
		Use:   "update [flags] <script-id>",
		Short: "Update an existing script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			scriptID := args[0]
			var content io.Reader
			if f.sourcePath != "" {
				c, err := os.Open(f.sourcePath)
				if err != nil {
					return err
				}
				defer c.Close()
				content = c
			}
			if err := script_update.NewUseCase(svc.ScriptService).Execute(ctx, scriptID, f.name, content); err != nil {
				return fmt.Errorf("updating script metadata: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the script")
	cmd.Flags().StringVar(&f.sourcePath, "source", "", "Script source file path")
	return cmd, nil
}
