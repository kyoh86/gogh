package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook/add"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name          string
		triggerEvent  string
		repoPattern   string
		operationType string
		operationID   string
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new hook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			opts := add.Options{
				Name:          f.name,
				TriggerEvent:  f.triggerEvent,
				RepoPattern:   f.repoPattern,
				OperationType: f.operationType,
				OperationID:   f.operationID,
			}
			id, err := add.NewUseCase(svc.HookService).Execute(ctx, opts)
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
	if err := enumFlag(cmd, &f.operationType, "operation-type", "", "Operation type", "overlay", "script"); err != nil {
		return nil, fmt.Errorf("registering operation-type flag: %w", err)
	}
	if err := cmd.MarkFlagRequired("operation-type"); err != nil {
		return nil, fmt.Errorf("marking operation-type flag required: %w", err)
	}
	cmd.Flags().StringVar(&f.operationID, "operation-id", "", "Operation resource ID")
	if err := cmd.MarkFlagRequired("operation-id"); err != nil {
		return nil, fmt.Errorf("marking operation-id flag required: %w", err)
	}
	return cmd, nil
}
