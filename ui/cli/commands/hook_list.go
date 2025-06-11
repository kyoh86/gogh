package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered hooks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			for h, err := range hook_list.NewUseCase(svc.HookService).Execute(ctx) {
				if err != nil {
					return err
				}
				fmt.Printf("* [%s] %s (%s)\n", h.ID, h.Name, h.Target)
			}
			return nil
		},
	}
	return cmd, nil
}
