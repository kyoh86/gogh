package commands

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v4/app/hook_edit"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookEditCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "edit [flags] <hook-id>",
		Short: "Edit an existing hook (edit Lua script with $EDITOR)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			// Extract the script to a temporary file
			tmpFile, err := os.CreateTemp("", "gogh_hook_edit_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile.Name())

			if err := hook_edit.NewUseCase(svc.HookService).ExtractScript(ctx, hookID, tmpFile); err != nil {
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
			return hook_edit.NewUseCase(svc.HookService).UpdateScript(ctx, hookID, edited)
		},
	}
	return cmd, nil
}
