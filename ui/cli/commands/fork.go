package commands

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/fork"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewForkCommand(
	defaultNames repository.DefaultNameService,
	tokens auth.TokenService,
	defaults *config.FlagStore,
	hostingService hosting.HostingService,
) *cobra.Command {
	var f config.ForkFlags

	checkFlags := func(_ *cobra.Command, args []string) (*repository.Reference, *repository.ReferenceWithAlias, error) {
		if len(args) != 1 {
			return nil, nil, fmt.Errorf("invalid number of arguments")
		}
		parser := repository.NewReferenceParser(defaultNames.GetDefaultHostAndOwner())
		srcRef, err := parser.Parse(args[0])
		if err != nil {
			return nil, nil, err
		}
		if f.To == "" {
			owner, err := defaultNames.GetDefaultOwnerFor(srcRef.Host())
			if err != nil {
				return nil, nil, err
			}
			return srcRef, &repository.ReferenceWithAlias{
				Reference: repository.NewReference(srcRef.Host(), owner, srcRef.Name()),
			}, nil
		}
		toRef, err := parser.ParseWithAlias(f.To)
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

	useCase := fork.NewUseCase(hostingService)

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
				//TODO: flag:Clone Retry Limit
				//TODO: flag:Default Branch Only
			}
			if err := useCase.Execute(ctx, *ref, *toRef, opts); err != nil {
				log.FromContext(ctx).WithError(err).Error("failed to fork the repository")
				return nil
			}
			return nil
		},
	}
	f.To = defaults.Fork.To
	// TODO: flag:Clone Retry Limit
	// TODO: flag:Default Branch Only
	cmd.Flags().
		StringVarP(
			&f.To,
			"to",
			"",
			"",
			strings.Join([]string{
				"Fork to the specified repository.",
				"It accepts a notation like 'OWNER/NAME' or 'OWNER/NAME=ALIAS'.",
				"If not specified, it will be forked to the default owner and same name as the original repository.",
				"If the alias is specified, it will be set as the local repository name.",
			}, " "),
		)
	return cmd
}
