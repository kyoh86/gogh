package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewForkCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.ForkFlags

	useCase := fork.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.OverlayService, svc.DefaultNameService, svc.ReferenceParser, svc.GitService)

	cmd := &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			opts := fork.Options{
				TryCloneOptions: try_clone.Options{
					Timeout: f.CloneRetryTimeout,
					Notify:  try_clone.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				},
				HostingOptions: fork.HostingOptions{
					DefaultBranchOnly: f.DefaultBranchOnly,
				},
				Target: f.To,
			}
			if err := useCase.Execute(ctx, refs[0], opts); err != nil {
				return fmt.Errorf("forking the repository: %w", err)
			}

			overlayFindUseCase := overlay_find.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayService,
			)
			overlayApplyUseCase := overlay_apply.NewUseCase()
			for overlay, err := range overlayFindUseCase.Execute(ctx, refs[0]) {
				if err != nil {
					return fmt.Errorf("finding overlays for %s: %w", refs[0], err)
				}
				var selected string
				if err := huh.NewForm(huh.NewGroup(
					huh.NewSelect[string]().
						Title("A repository to delete").
						Options(huh.Option[string]{
							Key:   "y",
							Value: "Yes",
						}, huh.Option[string]{
							Key:   "n",
							Value: "No",
						}, huh.Option[string]{
							Key:   "q",
							Value: "Quit",
						}).
						Value(&selected),
				)).Run(); err != nil {
					return err
				}
				switch selected {
				case "Yes", "y":
					if err := overlayApplyUseCase.Execute(ctx, overlay.Location.FullPath(), overlay.RelativePath, overlay.Content); err != nil {
						return err
					}
				case "No", "n":
					logger.Infof("Skipped applying overlay for %s", refs[0])
				case "Quit", "q":
					logger.Info("Quit applying overlays")
					return nil
				}
			}
			logger.Infof("Applied overlay for %s", refs[0])
			return nil
		},
	}
	cmd.Flags().StringVarP(
		&f.To,
		"to",
		"",
		svc.Flags.Fork.To,
		strings.Join([]string{
			"Fork to the specified repository.",
			"It accepts a notation like 'OWNER/NAME' or 'OWNER/NAME=ALIAS'.",
			"If not specified, it will be forked to the default owner and same name as the original repository.",
			"If the alias is specified, it will be set as the local repository name.",
		}, " "),
	)
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Fork.CloneRetryLimit, "")
	flags.BoolVarP(cmd, &f.DefaultBranchOnly, "default-branch-only", "", svc.Flags.Fork.DefaultBranchOnly, "Only fork the default branch")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Fork.CloneRetryTimeout, "Timeout for each clone attempt")
	return cmd, nil
}
