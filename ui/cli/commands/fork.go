package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewForkCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.ForkFlags

	cmd := &cobra.Command{
		Use:   "fork [flags] [<host>/]<owner>/<name>",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		Example: `  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host set by "config set-default-host"`,
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
			if err := fork.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.OverlayStore, svc.DefaultNameService, svc.ReferenceParser, svc.GitService).Execute(ctx, refs[0], opts); err != nil {
				return fmt.Errorf("forking the repository: %w", err)
			}

			var applied bool
			useCase := overlay_apply.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
				svc.OverlayStore,
			)
			if err := view.ProcessWithConfirmation(
				ctx,
				typ.FilterE(overlay_find.NewUseCase(
					svc.ReferenceParser,
					svc.OverlayStore,
				).Execute(ctx, refs[0]), func(ov *overlay_find.Overlay) (bool, error) {
					return !ov.ForInit, nil
				}),
				func(ov *overlay_find.Overlay) string {
					return fmt.Sprintf("Apply overlay for %s (%s)", refs[0], ov.RelativePath)
				},
				func(ov *overlay_find.Overlay) error {
					if err := useCase.Execute(ctx, refs[0], ov.RepoPattern, ov.ForInit, ov.RelativePath); err != nil {
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
				logger.Infof("Applied overlay for %s", refs[0])
			}
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
			"It accepts a notation like '<owner>/<name>' or '<owner>/<name>=<alias>'.",
			"If not specified, it will be forked to the default owner and same name as the original repository.",
			"If the alias is specified, it will be set as the local repository name",
		}, " "),
	)
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Fork.CloneRetryLimit, "The number of retries to clone a repository")
	flags.BoolVarP(cmd, &f.DefaultBranchOnly, "default-branch-only", "", svc.Flags.Fork.DefaultBranchOnly, "Only fork the default branch")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Fork.CloneRetryTimeout, "Timeout for each clone attempt")
	return cmd, nil
}
