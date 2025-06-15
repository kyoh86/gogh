package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v4/app/hook_describe"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookDescribeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "describe [flags] <hook-id>",
		Short: "Describe a hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			hook, code, err := hook_describe.NewUseCase(svc.HookService).Execute(ctx, hookID)
			if err != nil {
				return fmt.Errorf("describe hook: %w", err)
			}

			fmt.Println(strings.Repeat("-", 80))
			fmt.Printf("-- File Name: %s\n", hook.ScriptPath)
			fmt.Printf("-- Hook ID: %s\n", hook.ID)
			fmt.Printf("-- Hook Name: %s\n", hook.Name)
			if hook.Target.UseCase != "" {
				fmt.Printf("-- Target Use Case: %s\n", hook.Target.UseCase)
			}
			if hook.Target.Event != "" {
				fmt.Printf("-- Target Event: %s\n", hook.Target.Event)
			}
			if hook.Target.RepoPattern != "" {
				fmt.Printf("-- Target Repository Pattern: %s\n", hook.Target.RepoPattern)
			}
			fmt.Println(strings.Repeat("-", 80))
			if _, err := cmd.OutOrStdout().Write(code); err != nil {
				return fmt.Errorf("write hook code: %w", err)
			}
			return nil
		},
	}
	return cmd, nil
}
