package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/extra/apply"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	usecase := apply.NewUsecase(
		svc.ExtraService,
		svc.OverlayService,
		svc.WorkspaceService,
		svc.FinderService,
		svc.ReferenceParser,
	)

	var opts apply.Options

	cmd := &cobra.Command{
		Use:   "apply <name>",
		Short: "Apply a named extra to a repository",
		Long: `Apply a named extra template to a repository.

This applies all overlays in the named extra to the target repository.
By default, it applies to the current directory's repository.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			return usecase.Execute(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.TargetRepo, "target", "t", "", "Target repository (default: current directory)")

	return cmd, nil
}
