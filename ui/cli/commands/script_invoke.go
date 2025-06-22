package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/script_invoke"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptInvokeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "invoke [flags] <script-id> [[<host>/]<owner>/]<name>",
		Short: "Invoke an script in a repository",
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
			scriptID := args[0]
			ref := args[1]
			scriptInvokeUseCase := script_invoke.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ScriptService,
				svc.ReferenceParser,
			)
			if err := scriptInvokeUseCase.Execute(ctx, ref, scriptID, map[string]any{}); err != nil {
				return err
			}
			logger.Infof("Applied script %s to %s", scriptID, ref)
			return nil
		},
	}
	return cmd, nil
}
