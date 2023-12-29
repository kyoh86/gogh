package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

type forkFlagsStruct struct {
	Own bool `yaml:"own,omitempty"`
}

var (
	forkFlags forkFlagsStruct

	forkCommand = &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, specs []string) error {
			ctx := cmd.Context()
			parser := gogh.NewSpecParser(tokens.GetDefaultKey())
			spec, err := parser.Parse(specs[0])
			if err != nil {
				return err
			}
			_, token := tokens.Hosts.Get(spec.Host()).GetDefaultToken()
			adaptor, err := github.NewAdaptor(ctx, spec.Host(), token)
			if err != nil {
				return err
			}
			remote := gogh.NewRemoteController(adaptor)
			forked, err := remote.Fork(ctx, spec.Owner(), spec.Name(), nil)
			if err != nil {
				return err
			}

			root := defaultRoot()
			local := gogh.NewLocalController(root)

			localSpec := spec
			var opt *gogh.LocalCloneOption
			if forkFlags.Own {
				opt = &gogh.LocalCloneOption{Alias: &forked.Spec}
				localSpec = forked.Spec
			}
			log.FromContext(ctx).Infof("git clone %q", spec.URL())
			if _, err := local.Clone(ctx, spec, token, opt); err != nil {
				return fmt.Errorf("cloning the repository %q: %w", spec, err)
			}
			return local.SetRemoteSpecs(ctx, localSpec, map[string][]gogh.Spec{
				git.DefaultRemoteName: {forked.Spec},
				"upstream":            {spec},
			})
		},
	}
)

func init() {
	forkCommand.Flags().
		BoolVarP(&forkFlags.Own, "own", "", false, "Clones the forked repo to local as my-own repo")
	facadeCommand.AddCommand(forkCommand)
}
