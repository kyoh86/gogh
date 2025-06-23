package commands

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v4/app/extra_list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var (
		asJSON    bool
		extraType string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List extras",
		Long: `List all extras.

By default, lists all extras in one-line format.
Use --type to filter by extra type (auto, named).
Use --json to output in JSON format.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			useCase := extra_list.NewUseCase(svc.ExtraService, os.Stdout)
			return useCase.Execute(cmd.Context(), asJSON, extraType)
		},
	}

	cmd.Flags().StringVarP(&extraType, "type", "t", "all", "Filter by type (all, auto, named)")
	cmd.Flags().BoolVarP(&asJSON, "json", "j", false, "Output in JSON format")

	return cmd, nil
}
