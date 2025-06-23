package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name         string
		triggerEvent string
		repoPattern  string
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new hook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			opts := hook_add.Options{
				Name:         f.name,
				TriggerEvent: f.triggerEvent,
				RepoPattern:  f.repoPattern,
			}
			id, err := hook_add.NewUseCase(svc.HookService).Execute(ctx, opts)
			if err != nil {
				return fmt.Errorf("add hook: %w", err)
			}
			fmt.Printf("Hook added: %s\n", id)
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the hook")

	if err := enumFlag(cmd, &f.triggerEvent, "trigger-event", "", "event that triggers the hook", "", "post-clone", "post-fork", "post-create"); err != nil {
		return nil, fmt.Errorf("registering event flag: %w", err)
	}

	cmd.Flags().StringVar(&f.repoPattern, "repo-pattern", "", "Repository pattern")
	return cmd, nil
}
