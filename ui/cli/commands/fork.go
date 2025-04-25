package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
)

func NewForkCommand(conf *config.Config, tokens *config.TokenManager, defaults *config.Flags) *cobra.Command {
	var f config.ForkFlags
	cmd := &cobra.Command{
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
			_, token, err := tokens.GetDefaultTokenFor(spec.Host())
			if err != nil {
				return err
			}
			adaptor, err := github.NewAdaptor(ctx, spec.Host(), &token)
			if err != nil {
				return err
			}
			remote := gogh.NewRemoteController(adaptor)
			forked, err := remote.Fork(ctx, spec.Owner(), spec.Name(), nil)
			if err != nil {
				return err
			}

			root := conf.DefaultRoot()
			local := gogh.NewLocalController(root)

			localSpec := spec
			var opt *gogh.LocalCloneOption
			if f.Own {
				opt = &gogh.LocalCloneOption{Alias: &forked.Spec}
				localSpec = forked.Spec
			}
			log.FromContext(ctx).Infof("git clone %q", spec.URL())
			accessToken, err := adaptor.GetAccessToken()
			if err != nil {
				log.FromContext(ctx).WithField("error", err).Error("failed to get access token")
				return nil
			}
			if _, err := local.Clone(ctx, spec, accessToken, opt); err != nil {
				return fmt.Errorf("cloning the repository %q: %w", spec, err)
			}
			return local.SetRemoteSpecs(ctx, localSpec, map[string][]gogh.Spec{
				git.DefaultRemoteName: {forked.Spec},
				"upstream":            {spec},
			})
		},
	}
	f.Own = defaults.Fork.Own
	cmd.Flags().
		BoolVarP(&f.Own, "own", "", false, "Clones the forked repo to local as my-own repo")
	return cmd
}
