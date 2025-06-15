package commands

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewOverlayApplyCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		allRepos    bool
		repoPattern string
		forInit     bool
	}

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		listOpts := list.Options{
			Primary: false,
			ListOptions: list.ListOptions{
				Limit: 0,
			},
		}
		if f.allRepos {
			var ret []string
			for repo, err := range list.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx, listOpts) {
				if err != nil {
					return nil, fmt.Errorf("listing up repositories: %w", err)
				}
				ret = append(ret, repo.Path())
			}
			return ret, nil
		}
		var opts []huh.Option[string]
		for repo, err := range list.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx, listOpts) {
			if err != nil {
				return nil, fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Path(),
				Value: repo.Path(),
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
		Example: `  If you want to apply overlays to all repositories in the workspace,
  use --all-repositories flag.

  It accepts a short notation for each repository
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

			overlayApplyUseCase := overlay_apply.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayService,
			)
			var listup func(ctx context.Context, ref string) iter.Seq2[*overlay_find.Overlay, error]
			var filter func(ov *overlay_find.Overlay) (bool, error)
			if f.repoPattern == "" {
				listup = overlay_find.NewUseCase(svc.ReferenceParser, svc.OverlayService).Execute
				filter = func(ov *overlay_find.Overlay) (bool, error) {
					return ov.ForInit == f.forInit, nil // Filter by `forInit` flag
				}
			} else {
				listup = func(ctx context.Context, ref string) iter.Seq2[*overlay_find.Overlay, error] {
					return overlay_list.NewUseCase(svc.OverlayService).Execute(ctx)
				}
				filter = func(ov *overlay_find.Overlay) (bool, error) {
					return ov.RepoPattern == f.repoPattern && ov.ForInit == f.forInit, nil // Filter by `forInit` flag
				}
			}
			for _, ref := range refs {
				var applied bool
				if err := view.ProcessWithConfirmation(
					ctx,
					typ.FilterE(listup(ctx, ref), filter),
					func(ov *overlay_find.Overlay) string {
						return fmt.Sprintf("Apply overlay for %s (%s)", ref, ov.RelativePath)
					},
					func(ov *overlay_find.Overlay) error {
						if err := overlayApplyUseCase.Execute(ctx, ref, ov.RepoPattern, ov.ForInit, ov.RelativePath); err != nil {
							return err
						}
						applied = true
						return nil
					},
				); err != nil {
					if errors.Is(err, view.ErrQuit) {
						return nil
					}
					return err
				}
				if applied {
					logger.Infof("Applied overlay for %s", ref)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.allRepos, "all-repositories", "", false, "Apply overlays to all repositories in the workspace")
	cmd.Flags().StringVarP(&f.repoPattern, "repo-pattern", "", "", "Force apply overlays having this pattern, ignoring automatic repository name matching (useful for applying specific overlays or templates that would not normally match)")
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Apply overlays only for `gogh create` command (useful for templates)")
	return cmd, nil
}
