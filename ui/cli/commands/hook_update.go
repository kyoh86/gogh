package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook/update"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookUpdateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		name          string
		repoPattern   string
		triggerEvent  string
		operationType string
		operationID   string
	}
	cmd := &cobra.Command{
		Use:   "update [flags] <hook-id>",
		Short: "Update an existing hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			opts := update.Options{
				Name:          f.name,
				RepoPattern:   f.repoPattern,
				TriggerEvent:  f.triggerEvent,
				OperationType: f.operationType,
				OperationID:   f.operationID,
			}
			if err := update.NewUsecase(svc.HookService).Execute(ctx, hookID, opts); err != nil {
				return fmt.Errorf("updating hook metadata: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&f.name, "name", "", "Name of the hook")
	cmd.Flags().StringVar(&f.repoPattern, "repo-pattern", "", "Repository pattern")
	if err := enumFlag(cmd, &f.triggerEvent, "trigger-event", "", "event to hook automatically", "post-clone", "post-fork", "post-create"); err != nil {
		return nil, fmt.Errorf("registering trigger-event flag: %w", err)
	}
	if err := enumFlag(cmd, &f.operationType, "operation-type", "", "Operation type", "overlay", "script"); err != nil {
		return nil, fmt.Errorf("registering operation-type flag: %w", err)
	}
	cmd.Flags().StringVar(&f.operationID, "operation-id", "", "Operation resource ID")
	return cmd, nil
}
