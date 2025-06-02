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
	reposUseCase := repos.NewUseCase(svc.HostingService)

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
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
		Use:   "apply [[<owner>/]<name>...]",
		Short: "Target overlays to a repository",
		Args:  cobra.ArbitraryArgs,
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
				svc.OverlayService,
			)
			overlayApplyUseCase := overlay_apply.NewUseCase(svc.OverlayService)
			for _, ref := range refs {
				if err := view.ProcessWithConfirmation(
					ctx,
					typ.Filter2(overlayFindUseCase.Execute(ctx, ref), func(entry *overlay_find.OverlayEntry) bool {
						return !entry.ForInit
					}),
					func(entry *overlay_find.OverlayEntry) string {
						return fmt.Sprintf("Apply overlay for %s (%s)", ref, entry.RelativePath)
					},
					func(entry *overlay_find.OverlayEntry) error {
						return overlayApplyUseCase.Execute(ctx, entry.Location.FullPath(), entry.Pattern, entry.ForInit, entry.RelativePath)
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

	return cmd, nil
}
