package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/hook/invoke"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookInvokeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "invoke [flags] <hook-id> [[<host>/]<owner>/]<name>",
		Short: "Manually invoke a hook for a repository",
		Args:  cobra.ExactArgs(2),
		Example: `  invoke <hook-id> github.com/owner/repo
  invoke <hook-id> owner/repo
  invoke <hook-id> repo
  invoke <hook-id> .  # Use current directory repository
  
  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hookID := args[0]
			repoRef := args[1]

			// Use current directory if reference is "."
			if repoRef == "." {
				repo, err := cwd.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
				if err != nil {
					return fmt.Errorf("finding repository from current directory: %w", err)
				}
				repoRef = repo.Ref().String()
			}

			return invoke.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.HookService,
				svc.OverlayService,
				svc.ScriptService,
				svc.ReferenceParser,
			).Invoke(ctx, hookID, repoRef)
		},
	}
	return cmd, nil
}
