package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/hook/describe"
	"github.com/kyoh86/gogh/v4/app/hook/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		json   bool
		source bool
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List registered hooks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var usecase interface {
				Execute(ctx context.Context, s describe.Hook) error
			}
			if f.json {
				usecase = describe.NewJSONUsecase(cmd.OutOrStdout())
			} else {
				usecase = describe.NewOnelineUsecase(cmd.OutOrStdout())
			}
			for s, err := range list.NewUsecase(svc.HookService).Execute(cmd.Context()) {
				if err != nil {
					return err
				}
				if s == nil {
					continue
				}
				if err := usecase.Execute(cmd.Context(), s); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.json, "json", "", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&f.source, "source", "", false, "Output with source code")
	return cmd, nil
}
