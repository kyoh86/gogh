package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook/remove"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "remove [flags] <hook-id>",
		Short: "Remove a registered hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			return remove.NewUseCase(svc.HookService).Execute(ctx, hookID)
		},
	}
	return cmd, nil
}
