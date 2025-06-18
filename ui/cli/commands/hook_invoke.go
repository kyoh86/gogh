package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook_invoke"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookInvokeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "invoke [flags] <hook-id> [[<host>/]<owner>/]<name>",
		Short: "Run a hook script forcely for a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			repoRef := args[1]
			return hook_invoke.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.HookService,
				svc.OverlayService,
				svc.ScriptService,
				svc.ReferenceParser,
			).Invoke(ctx, hookID, repoRef)
		},
	}
	return cmd, nil
}
