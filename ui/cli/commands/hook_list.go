package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook_list"
	"github.com/kyoh86/gogh/v4/app/hook_show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json bool
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered hooks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			showUseCase := hook_show.NewUseCase(cmd.OutOrStdout(), f.json)
			for h, err := range hook_list.NewUseCase(svc.HookService).Execute(ctx) {
				if err != nil {
					return err
				}
				if err := showUseCase.Execute(ctx, h); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "j", false, "Output in JSON format")
	return cmd, nil
}
