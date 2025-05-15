package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/fork"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewForkCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	var f config.ForkFlags

	useCase := fork.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.DefaultNameService, svc.ReferenceParser, svc.GitService)

	cmd := &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			opts := fork.Options{
				TryCloneNotify: service.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				HostingOptions: fork.HostingOptions{
					DefaultBranchOnly: f.DefaultBranchOnly,
				},
				Target: f.To,
			}
			if err := useCase.Execute(ctx, refs[0], opts); err != nil {
				return fmt.Errorf("failed to fork the repository: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().
		StringVarP(
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
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Create.CloneRetryLimit, "")
	cmd.Flags().
		BoolVarP(&f.DefaultBranchOnly, "default-branch-only", "", false, "Only fork the default branch")
	return cmd
}
