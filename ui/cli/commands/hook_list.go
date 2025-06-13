package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/hook"
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
			write := func(h *hook.Hook) {
				fmt.Printf("* [%s] %s (%s)\n", h.ID, h.Name, h.Target)
			}
			if f.json {
				enc := json.NewEncoder(cmd.OutOrStdout())
				write = func(h *hook.Hook) {
					if err := enc.Encode(h); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error encoding hook %s: %v\n", h.ID, err)
					}
				}
			}
			for h, err := range hook_list.NewUseCase(svc.HookService).Execute(ctx) {
				if err != nil {
					return err
				}
				write(h)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "j", false, "Output in JSON format")
	return cmd, nil
}
