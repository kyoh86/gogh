package commands

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v4/app/extra/show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "show <id-or-name>",
		Short: "Show details of an extra",
		Long: `Show detailed information about an extra.

You can specify either an extra ID or name (for named extras).
Use --json to output in JSON format.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			usecase := show.NewUsecase(svc.ExtraService, os.Stdout)
			return usecase.Execute(cmd.Context(), args[0], asJSON)
		},
	}

	cmd.Flags().BoolVarP(&asJSON, "json", "j", false, "Output in JSON format")

	return cmd, nil
}
