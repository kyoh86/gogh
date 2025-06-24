package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook/show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json   bool
		source bool
	}
	cmd := &cobra.Command{
		Use:   "show [flags] <hook-id>",
		Short: "Show a hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			return show.NewUseCase(svc.HookService, cmd.OutOrStdout()).Execute(ctx, hookID, f.json)
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "", false, "Output in JSON format")
	return cmd, nil
}
