package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/extra_save"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraSaveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "save <repository>",
		Short: "Save excluded files as auto-apply extra",
		Long: `Save files that are excluded by .gitignore as auto-apply extra.
These extra will be automatically applied when the repository is cloned.`,
		Args: cobra.ExactArgs(1),
		Example: `  save github.com/kyoh86/example
  save .  # Save from current directory repository

  It accepts a short notation for the repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			repoStr := args[0]

			// Use current directory if reference is "."
			if repoStr == "." {
				repo, err := cwd.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
				if err != nil {
					return fmt.Errorf("finding repository from current directory: %w", err)
				}
				repoStr = repo.Ref().String()
			}

			uc := extra_save.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.GitService,
				svc.OverlayService,
				svc.HookService,
				svc.ExtraService,
				svc.ReferenceParser,
			)

			if err := uc.Execute(ctx, repoStr); err != nil {
				return err
			}

			logger.Infof("Saved auto-apply extra for %s", repoStr)
			return nil
		},
	}
	return cmd, nil
}
