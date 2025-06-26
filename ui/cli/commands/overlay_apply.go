package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/app/overlay/apply"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		allRepositories bool
		patterns        []string
	}
	cmd := &cobra.Command{
		Use:   "apply [flags] <overlay-id> [[<host>/]<owner>/]<name>",
		Short: "Apply an overlay to a repository",
		Args:  cobra.MinimumNArgs(1),
		Example: `  invoke [flags] <overlay-id> [[[<host>/]<owner>/]<name>...]
  invoke [flags] <overlay-id> --all
  invoke [flags] <overlay-id> --pattern <pattern> [--pattern <pattern>]...

  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
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
			refs := args[1:]
			overlayApplyUsecase := apply.NewUsecase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayService,
			)
			if f.allRepositories || len(f.patterns) > 0 {
				if len(refs) > 0 {
					return errors.New("cannot specify repositories when --all or --pattern flag is set")
				}

				// If --all flag is set, apply the script to all repositories in the workspace
				for repo, err := range list.NewUsecase(
					svc.WorkspaceService,
					svc.FinderService,
				).Execute(ctx, list.Options{ListOptions: list.ListOptions{
					Limit:    0,
					Patterns: f.patterns,
				}}) {
					if err != nil {
						return fmt.Errorf("listing repositories: %w", err)
					}
					refs = append(refs, repo.Ref().String())
				}
			}
			for _, ref := range refs {
				// Use current directory if reference is "."
				if ref == "." {
					repo, err := cwd.NewUsecase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
					if err != nil {
						return fmt.Errorf("finding repository from current directory: %w", err)
					}
					ref = repo.Ref().String()
				}

				if err := overlayApplyUsecase.Execute(ctx, ref, overlayID); err != nil {
					return err
				}
				logger.Infof("Applied overlay %s to %s", overlayID, ref)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.allRepositories, "all", "", false, "Apply to all repositories in the workspace")
	cmd.Flags().StringSliceVarP(&f.patterns, "pattern", "p", nil, "Patterns for selecting repositories")
	return cmd, nil
}
