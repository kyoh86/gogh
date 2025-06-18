package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "apply [flags] <overlay-id> [[<host>/]<owner>/]<name>",
		Short: "Apply an overlay to a repository",
		Args:  cobra.ExactArgs(2),
		Example: `  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			overlayID := args[0]
			ref := args[1]
			overlayApplyUseCase := overlay_apply.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayService,
			)
			if err := overlayApplyUseCase.Execute(ctx, ref, overlayID); err != nil {
				return err
			}
			logger.Infof("Applied overlay %s to %s", overlayID, ref)
			return nil
		},
	}
	return cmd, nil
}
