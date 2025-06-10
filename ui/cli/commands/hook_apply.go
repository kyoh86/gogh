package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook_apply"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "apply [flags] <hook-id> [[<host>/]<owner>/]<name>",
		Short: "Run a hook script forcely for a repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			repoRef := args[1]
			return hook_apply.NewUseCase(
				svc.HookService,
				svc.ReferenceParser,
				svc.WorkspaceService,
				svc.FinderService,
			).Execute(ctx, hookID, repoRef, map[string]any{
				"use_case": "hook-apply",
				"event":    "",
			})
		},
	}
	return cmd, nil
}
