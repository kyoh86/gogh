package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/fork"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
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
			if err := fork.
				NewUseCase(
					svc.HostingService,
					svc.WorkspaceService,
					svc.FinderService,
					svc.OverlayService,
					svc.ScriptService,
					svc.HookService,
					svc.DefaultNameService,
					svc.ReferenceParser,
					svc.GitService,
				).
				Execute(ctx, refs[0], opts); err != nil {
				return fmt.Errorf("forking the repository: %w", err)
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
	cmd.Flags().BoolVarP(&f.DefaultBranchOnly, "default-branch-only", "", svc.Flags.Fork.DefaultBranchOnly, "Only fork the default branch")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Fork.CloneRetryTimeout, "Timeout for each clone attempt")
	return cmd, nil
}
