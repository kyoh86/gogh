package commands

import (
	"context"
	"os"

	scriptedit "github.com/kyoh86/gogh/v4/app/script/edit"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/completion"
	"github.com/spf13/cobra"
)

func NewScriptEditCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "edit [flags] <script-id>",
		Short: "Edit an existing script (with $EDITOR)",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completion.Scripts(cmd.Context(), svc, toComplete)
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			scriptID := args[0]
			// Extract the script to a temporary file
			tmpFile, err := os.CreateTemp("", "gogh_script_edit_*.lua")
			if err != nil {
				return err
			}
			defer os.Remove(tmpFile.Name())

			if err := scriptedit.NewUsecase(svc.ScriptService).ExtractScript(ctx, scriptID, tmpFile); err != nil {
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
			return scriptedit.NewUsecase(svc.ScriptService).UpdateScript(ctx, scriptID, edited)
		},
	}
	return cmd, nil
}
