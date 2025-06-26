package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/script/show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json   bool
		source bool
	}
	cmd := &cobra.Command{
		Use:   "show [flags] <script-id>",
		Short: "Show a script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			scriptID := args[0]
			return show.NewUsecase(svc.ScriptService, cmd.OutOrStdout()).Execute(ctx, scriptID, f.json, f.source)
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&f.source, "source", "", false, "Output with source code")
	return cmd, nil
}
