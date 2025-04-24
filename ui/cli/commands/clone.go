package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/cmdutil"
	"github.com/kyoh86/gogh/v3/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCloneCommand(conf *config.Config, tokens *config.TokenManager) *cobra.Command {
	var f struct {
		dryrun bool
	}

	c := &cobra.Command{
		Use:     "clone [flags] [[OWNER/]NAME[=ALIAS]]...",
		Aliases: []string{"get"},
		Short:   "Clone repositories to local",
		Example: `  It accepts a shortly notation for a repository
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

		RunE: func(cmd *cobra.Command, specs []string) error {
			ctx := cmd.Context()
			if len(specs) == 0 {
				entries := tokens.Entries()
				var options []huh.Option[string]
				for _, entry := range entries {
					adaptor, err := github.NewAdaptor(ctx, entry.Host, &entry.Token)
					if err != nil {
						return err
					}
					remote := gogh.NewRemoteController(adaptor)
					founds, err := remote.List(ctx, nil)
					if err != nil {
						return err
					}
					for _, s := range founds {
						options = append(options, huh.Option[string]{
							Key:   s.Spec.String(),
							Value: s.Spec.String(),
						})
					}
				}
				if err := huh.NewForm(huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("A repository to clone").
						Options(options...).
						Value(&specs),
				)).Run(); err != nil {
					return err
				}
			}
			return cloneAll(ctx, conf, tokens, specs, f.dryrun)
		},
	}

	c.Flags().
		BoolVarP(&f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return c
}

func cloneAll(ctx context.Context, conf *config.Config, tokens *config.TokenManager, specs []string, dryrun bool) error {
	parser := gogh.NewSpecParser(tokens.GetDefaultKey())
	if dryrun {
		for _, s := range specs {
			spec, alias, err := parser.ParseWithAlias(s)
			if err != nil {
				return err
			}

			if alias == nil {
				log.FromContext(ctx).Infof("git clone %q", spec.URL())
			} else {
				log.FromContext(ctx).Infof("git clone %q into %q", spec.URL(), alias.String())
			}
		}
		return nil
	}

	local := gogh.NewLocalController(conf.DefaultRoot())
	if len(specs) == 1 {
		return cloneOneFunc(ctx, tokens, local, parser, specs[0])()
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, s := range specs {
		eg.Go(cloneOneFunc(ctx, tokens, local, parser, s))
	}
	return eg.Wait()
}

func cloneOneFunc(
	ctx context.Context,
	tokens *config.TokenManager,
	local *gogh.LocalController,
	parser gogh.SpecParser,
	s string,
) func() error {
	return func() error {
		spec, alias, err := parser.ParseWithAlias(s)
		if err != nil {
			return err
		}

		adaptor, remote, err := cmdutil.RemoteControllerFor(ctx, *tokens, spec)
		if err != nil {
			return err
		}
		repo, err := remote.Get(ctx, spec.Owner(), spec.Name(), nil)
		if err != nil {
			return err
		}

		l := log.FromContext(ctx).WithFields(log.Fields{
			"spec": spec,
		})
		l.Info("cloning")
		accessToken, err := adaptor.GetAccessToken()
		if err != nil {
			l.WithField("error", err).Error("failed to get access token")
			return nil
		}
		if _, err = local.Clone(ctx, spec, accessToken, &gogh.LocalCloneOption{Alias: alias}); err != nil {
			l.WithField("error", err).Error("failed to get repository")
			return nil
		}
		if repo.Parent != nil && repo.Parent.String() != spec.String() {
			l.Debug("set remote specs")
			localSpec := spec
			if alias != nil {
				localSpec = *alias
			}
			if err := local.SetRemoteSpecs(ctx, localSpec, map[string][]gogh.Spec{
				git.DefaultRemoteName: {spec},
				"upstream":            {*repo.Parent},
			}); err != nil {
				return err
			}
		}
		l.Info("finished")
		return nil
	}
}
