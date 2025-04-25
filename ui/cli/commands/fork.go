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

func NewForkCommand(conf *config.ConfigStore, tokens *config.TokenStore, defaults *config.FlagStore) *cobra.Command {
	var f config.ForkFlags
	cmd := &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			parser := gogh.NewRepoRefParser(tokens.GetDefaultKey())
			ref, err := parser.Parse(refs[0])
			if err != nil {
				return err
			}
			_, token, err := tokens.GetDefaultTokenFor(ref.Host())
			if err != nil {
				return err
			}
			adaptor, err := github.NewAdaptor(ctx, ref.Host(), &token)
			if err != nil {
				return err
			}
			remote := gogh.NewRemoteController(adaptor)
			forked, err := remote.Fork(ctx, ref.Owner(), ref.Name(), nil)
			if err != nil {
				return err
			}

			root := conf.DefaultRoot()
			local := gogh.NewLocalController(root)

			localRef := ref
			var opt *gogh.LocalCloneOption
			if f.Own {
				opt = &gogh.LocalCloneOption{Alias: &forked.Ref}
				localRef = forked.Ref
			}
			log.FromContext(ctx).Infof("git clone %q", ref.URL())
			accessToken, err := adaptor.GetAccessToken()
			if err != nil {
				log.FromContext(ctx).WithField("error", err).Error("failed to get access token")
				return nil
			}
			if _, err := local.Clone(ctx, ref, accessToken, opt); err != nil {
				return fmt.Errorf("cloning the remote repository %q: %w", ref, err)
			}
			return local.SetRemoteRefs(ctx, localRef, map[string][]gogh.RepoRef{
				git.DefaultRemoteName: {forked.Ref},
				"upstream":            {ref},
			})
		},
	}
	f.Own = defaults.Fork.Own
	cmd.Flags().
		BoolVarP(&f.Own, "own", "", false, "Clones the forked remote repo to local as my-own repo")
	return cmd
}
