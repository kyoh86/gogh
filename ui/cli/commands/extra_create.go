package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/extra_create"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	useCase := extra_create.NewUseCase(
		svc.WorkspaceService,
		svc.FinderService,
		svc.ExtraService,
		svc.OverlayService,
		svc.ReferenceParser,
	)

	var opts extra_create.Options

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a named extra template",
		Long: `Create a named extra template from overlays.

This creates a reusable template that can be applied to any repository later.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return useCase.Execute(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.SourceRepo, "source", "s", "", "Source repository")
	cmd.Flags().StringSliceVarP(&opts.OverlayNames, "overlay", "o", nil, "Overlay names to include in the extra")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("overlay")

	return cmd, nil
}
