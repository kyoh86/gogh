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

	checkFlags := func(ctx context.Context, args []string) (string, error) {
		if len(args) != 0 {
			return args[0], nil
		}
		var opts []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
			if err != nil {
				return "", fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("A repository to delete").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return "", err
		}
		return selected, nil
	}

	cmd := &cobra.Command{
		Use:   "apply [[<owner>/]<name>]",
		Short: "Target overlays to a repository",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			selected, err := checkFlags(ctx, refs)
			if err != nil {
				return err
			}

			overlayFindUseCase := overlay_find.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayService,
			)
			overlayApplyUseCase := overlay_apply.NewUseCase()
			if err := view.ProcessWithConfirmation(
				ctx,
				typ.Filter2(overlayFindUseCase.Execute(ctx, selected), func(overlay *overlay_find.Overlay) bool {
					return !overlay.ForInit
				}),
				func(overlay *overlay_find.Overlay) string {
					return fmt.Sprintf("Apply overlay for %s (%s)", selected, overlay.RelativePath)
				},
				func(overlay *overlay_find.Overlay) error {
					return overlayApplyUseCase.Execute(ctx, overlay.Location.FullPath(), overlay.RelativePath, overlay.Content)
				},
			); err != nil {
				if errors.Is(err, view.ErrQuit) {
					return nil
				}
				return err
			}

			logger.Infof("Applied overlay for %s", selected)
			return nil
		},
	}

	return cmd, nil
}
