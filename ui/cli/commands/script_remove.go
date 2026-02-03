package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/script/remove"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/completion"
	"github.com/spf13/cobra"
)

func NewScriptRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "remove [flags] <script-id>",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove a script",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completion.Scripts(cmd.Context(), svc, toComplete)
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			scriptID := args[0]
			if err := remove.NewUsecase(svc.ScriptService).Execute(ctx, scriptID); err != nil {
				return err
			}

			logger.Infof("Removed script %s", scriptID)
			return nil
		},
	}
	return cmd, nil
}
