package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/extra/remove"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	useCase := remove.NewUseCase(
		svc.ExtraService,
		svc.ReferenceParser,
	)

	var opts remove.Options

	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an extra",
		Long: `Remove an extra.

You can remove by:
- ID: Use --id flag
- Name (for named extras): Use --name flag
- Repository (for auto extras): Use --repository flag`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return useCase.Execute(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.ID, "id", "i", "", "Extra ID to remove")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Named extra to remove")
	cmd.Flags().StringVarP(&opts.Repository, "repository", "r", "", "Repository whose auto extra to remove")
	cmd.MarkFlagsMutuallyExclusive("id", "name", "repository")

	return cmd, nil
}
