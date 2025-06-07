package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_apply"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/spf13/cobra"
)

func NewHookApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "apply <hook-id>",
		Short: "Run a hook script (for testing)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			var found *hook.Hook
			for h, err := range svc.HookService.ListHooks() {
				if err != nil {
					return err
				}
				if h.ID == hookID {
					found = h
					break
				}
			}
			if found == nil {
				return fmt.Errorf("hook not found: %s", hookID)
			}
			return hook_apply.NewUseCase(svc.HookService).Execute(ctx, *found, map[string]string{})
		},
	}
	return cmd, nil
}
