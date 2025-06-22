package commands

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v4/app/script_edit"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptEditCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "edit [flags] <script-id>",
		Short: "Edit an existing script (with $EDITOR)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			scriptID := args[0]
			// Extract the script to a temporary file
			tmpFile, err := os.CreateTemp("", "gogh_script_edit_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile.Name())

			if err := script_edit.NewUseCase(svc.ScriptService).ExtractScript(ctx, scriptID, tmpFile); err != nil {
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
			return script_edit.NewUseCase(svc.ScriptService).UpdateScript(ctx, scriptID, edited)
		},
	}
	return cmd, nil
}
