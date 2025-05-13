package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/fork"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewForkCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	var f config.ForkFlags

	checkFlags := func(_ *cobra.Command, args []string) (*repository.Reference, *repository.ReferenceWithAlias, error) {
		if len(args) != 1 {
			return nil, nil, fmt.Errorf("invalid number of arguments")
		}
		srcRef, err := svc.referenceParser.Parse(args[0])
		if err != nil {
			return nil, nil, err
		}
		if f.To == "" {
			owner, err := svc.defaultNameService.GetDefaultOwnerFor(srcRef.Host())
			if err != nil {
				return nil, nil, err
			}
			return srcRef, &repository.ReferenceWithAlias{
				Reference: repository.NewReference(srcRef.Host(), owner, srcRef.Name()),
			}, nil
		}
		toRef, err := svc.referenceParser.ParseWithAlias(f.To)
		if err != nil {
			return nil, nil, err
		}
		if toRef.Reference.Host() != srcRef.Host() {
			return nil, nil, fmt.Errorf("the host of the forked repository must be the same as the original repository")
		}
		if toRef.Reference.Owner() == "" {
			return nil, nil, fmt.Errorf("the owner of the forked repository must be specified")
		}
		return srcRef, toRef, nil
	}

	useCase := fork.NewUseCase(svc.hostingService, svc.workspaceService, svc.gitService)

	cmd := &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			ref, toRef, err := checkFlags(cmd, refs)
			if err != nil {
				return err
			}
			opts := fork.Options{
				TryCloneNotify: service.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				ForkRepositoryOptions: hosting.ForkRepositoryOptions{
					DefaultBranchOnly: f.DefaultBranchOnly,
				},
			}
			if err := useCase.Execute(ctx, *ref, *toRef, opts); err != nil {
				log.FromContext(ctx).WithError(err).Error("failed to fork the repository")
				return nil
			}
			return nil
		},
	}
	cmd.Flags().
		StringVarP(
			&f.To,
			"to",
			"",
			svc.flags.Fork.To,
			strings.Join([]string{
				"Fork to the specified repository.",
				"It accepts a notation like 'OWNER/NAME' or 'OWNER/NAME=ALIAS'.",
				"If not specified, it will be forked to the default owner and same name as the original repository.",
				"If the alias is specified, it will be set as the local repository name.",
			}, " "),
		)
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.flags.Create.CloneRetryLimit, "")
	cmd.Flags().
		BoolVarP(&f.DefaultBranchOnly, "default-branch-only", "", false, "Only fork the default branch")
	return cmd
}
