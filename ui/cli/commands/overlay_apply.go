package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewOverlayApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		repoPattern string
		forInit     bool
	}

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range repos.NewUseCase(svc.HostingService).Execute(ctx, repos.Options{}) {
			if err != nil {
				return nil, fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected []string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Repositories to apply overlays").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return nil, err
		}
		return selected, nil
	}

	cmd := &cobra.Command{
		Use:   "apply [flags] [[[<host>/]<owner>/]<name>...]",
		Short: "Apply overlays to specified repositories",
		Args:  cobra.ArbitraryArgs,
		Example: `  It accepts a short notation for a repository
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
			refs, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}

			overlayFindUseCase := overlay_find.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayStore,
			)
			overlayApplyUseCase := overlay_apply.NewUseCase(svc.OverlayStore)
			for _, ref := range refs {
				if err := view.ProcessWithConfirmation(
					ctx,
					typ.Filter2(overlayFindUseCase.Execute(ctx, ref), func(entry *overlay_find.OverlayEntry) bool {
						return f.forInit == entry.ForInit // Filter by `forInit` flag
					}),
					func(entry *overlay_find.OverlayEntry) string {
						return fmt.Sprintf("Apply overlay for %s (%s)", ref, entry.RelativePath)
					},
					func(entry *overlay_find.OverlayEntry) error {
						return overlayApplyUseCase.Execute(ctx, entry.Location.FullPath(), entry.RepoPattern, entry.ForInit, entry.RelativePath)
					},
				); err != nil {
					if errors.Is(err, view.ErrQuit) {
						return nil
					}
					return err
				}

				logger.Infof("Applied overlay for %s", ref)
			}
			return nil
		},
	}
	//TODO: Add flags for and --repo-pattern
	cmd.Flags().StringVarP(&f.repoPattern, "repo-pattern", "", "", "Force apply overlays having this pattern, ignoring automatic repository name matching (useful for applying specific overlays or templates that would not normally match)")
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Apply overlays only for `gogh create` command (useful for templates)")
	return cmd, nil
}
