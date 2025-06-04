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

			useCase := overlay_apply.NewUseCase(svc.OverlayStore)
			if err := view.ProcessWithConfirmation(
				ctx,
				typ.Filter2(overlay_find.NewUseCase(
					svc.WorkspaceService,
					svc.FinderService,
					svc.ReferenceParser,
					svc.OverlayStore,
				).Execute(ctx, refs[0]), func(entry *overlay_find.OverlayEntry) bool {
					return !entry.ForInit
				}),
				func(entry *overlay_find.OverlayEntry) string {
					return fmt.Sprintf("Apply overlay for %s (%s)", refs[0], entry.RelativePath)
				},
				func(entry *overlay_find.OverlayEntry) error {
					return useCase.Execute(ctx, entry.Location.FullPath(), entry.RepoPattern, entry.ForInit, entry.RelativePath)
				},
			); err != nil {
				if errors.Is(err, view.ErrQuit) {
					return nil
				}
				return err
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
