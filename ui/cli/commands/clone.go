package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/domain/local"
	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCloneCommand(conf *config.ConfigStore, defaultNames repository.DefaultNameService, tokens auth.TokenService) *cobra.Command {
	var f struct {
		dryrun bool
	}

	c := &cobra.Command{
		Use:     "clone [flags] [[OWNER/]NAME[=ALIAS]]...",
		Aliases: []string{"get"},
		Short:   "Clone remote repositories to local",
		Example: `  It accepts a shortly notation for a remote repository
  (for example, "github.com/kyoh86/example") like below.
    - "NAME": e.g. "example"; 
    - "OWNER/NAME": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default".

  It accepts an alias for each repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example"
  into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,

		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			if len(refs) == 0 {
				entries := tokens.Entries()
				var options []huh.Option[string]
				for _, entry := range entries {
					adaptor, err := github.NewAdaptor(ctx, entry.Host, &entry.Token)
					if err != nil {
						return err
					}
					ctrl := remote.NewController(adaptor)
					founds, err := ctrl.List(ctx, nil)
					if err != nil {
						return err
					}
					for _, s := range founds {
						options = append(options, huh.Option[string]{
							Key:   s.Ref.String(),
							Value: s.Ref.String(),
						})
					}
				}
				if err := huh.NewForm(huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("A remote repository to clone").
						Options(options...).
						Value(&refs),
				)).Run(); err != nil {
					return err
				}
			}
			return cloneAll(ctx, conf, defaultNames, tokens, refs, f.dryrun)
		},
	}

	c.Flags().
		BoolVarP(&f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return c
}

func cloneAll(ctx context.Context, conf *config.ConfigStore, defaultNames repository.DefaultNameService, tokens auth.TokenService, refs []string, dryrun bool) error {
	parser := reporef.NewRepoRefParser(defaultNames.GetDefaultHostAndOwner())
	if dryrun {
		for _, r := range refs {
			ref, alias, err := parser.ParseWithAlias(r)
			if err != nil {
				return err
			}

			if alias == nil {
				log.FromContext(ctx).Infof("git clone %q", ref.URL())
			} else {
				log.FromContext(ctx).Infof("git clone %q into %q", ref.URL(), alias.String())
			}
		}
		return nil
	}

	ctrl := local.NewController(conf.DefaultRoot())
	if len(refs) == 1 {
		return cloneOneFunc(ctx, tokens, ctrl, parser, refs[0])()
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, s := range refs {
		eg.Go(cloneOneFunc(ctx, tokens, ctrl, parser, s))
	}
	return eg.Wait()
}

func cloneOneFunc(
	ctx context.Context,
	tokens auth.TokenService,
	ctrl *local.Controller,
	parser reporef.RepoRefParser,
	s string,
) func() error {
	return func() error {
		ref, alias, err := parser.ParseWithAlias(s)
		if err != nil {
			return err
		}

		adaptor, remote, err := RemoteControllerFor(ctx, tokens, ref)
		if err != nil {
			return err
		}
		repo, err := remote.Get(ctx, ref.Owner(), ref.Name(), nil)
		if err != nil {
			return err
		}

		l := log.FromContext(ctx).WithFields(log.Fields{
			"ref": ref,
		})
		l.Info("cloning")
		accessToken, err := adaptor.GetAccessToken()
		if err != nil {
			l.WithField("error", err).Error("failed to get access token")
			return nil
		}
		if _, err = ctrl.Clone(ctx, ref, accessToken, &local.CloneOption{Alias: alias}); err != nil {
			l.WithField("error", err).Error("failed to clone a repository")
			return nil
		}
		if repo.Parent != nil && repo.Parent.String() != ref.String() {
			l.Debug("set remote refs")
			localRef := ref
			if alias != nil {
				localRef = *alias
			}
			if err := ctrl.SetRemoteRefs(ctx, localRef, map[string][]reporef.RepoRef{
				git.DefaultRemoteName: {ref},
				"upstream":            {*repo.Parent},
			}); err != nil {
				return err
			}
		}
		l.Info("finished")
		return nil
	}
}
