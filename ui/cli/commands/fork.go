package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/domain/local"
	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
)

func NewForkCommand(conf *config.ConfigStore, defaultNames repository.DefaultNameService, tokens auth.TokenService, defaults *config.FlagStore) *cobra.Command {
	var f config.ForkFlags
	cmd := &cobra.Command{
		Use:   "fork [flags] OWNER/NAME",
		Short: "Fork a repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			parser := reporef.NewRepoRefParser(defaultNames.GetDefaultHostAndOwner())
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
			remote := remote.NewController(adaptor)
			forked, err := remote.Fork(ctx, ref.Owner(), ref.Name(), nil)
			if err != nil {
				return err
			}

			root := conf.PrimaryRoot()
			ctrl := local.NewController(root)

			localRef := ref
			var opts *local.CloneOption
			if f.Own {
				opts = &local.CloneOption{Alias: &forked.Ref}
				localRef = forked.Ref
			}
			log.FromContext(ctx).Infof("git clone %q", ref.URL())
			accessToken, err := adaptor.GetAccessToken()
			if err != nil {
				log.FromContext(ctx).WithField("error", err).Error("failed to get access token")
				return nil
			}
			if _, err := ctrl.Clone(ctx, ref, accessToken, opts); err != nil {
				return fmt.Errorf("cloning the remote repository %q: %w", ref, err)
			}
			return ctrl.SetRemoteRefs(ctx, localRef, map[string][]reporef.RepoRef{
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
