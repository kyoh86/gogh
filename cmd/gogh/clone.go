package main

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var cloneFlags struct {
	dryrun bool
}

var cloneCommand = &cobra.Command{
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
			for _, entry := range tokens.Entries() {
				adaptor, err := github.NewAdaptor(ctx, entry.Host, entry.Token)
				if err != nil {
					return err
				}
				remote := gogh.NewRemoteController(adaptor)
				founds, err := remote.List(ctx, nil)
				if err != nil {
					return err
				}
				for _, s := range founds {
					specs = append(specs, s.Spec.String())
				}
			}
			if err := survey.AskOne(&survey.MultiSelect{
				Message: "A repository to clone",
				Options: specs,
			}, &specs); err != nil {
				return err
			}
		}
		return cloneAll(ctx, specs, cloneFlags.dryrun)
	},
}

func cloneAll(ctx context.Context, specs []string, dryrun bool) error {
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

	local := gogh.NewLocalController(defaultRoot())
	if len(specs) == 1 {
		return cloneOneFunc(ctx, local, parser, specs[0])()
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, s := range specs {
		eg.Go(cloneOneFunc(ctx, local, parser, s))
	}
	return eg.Wait()
}

func cloneOneFunc(
	ctx context.Context,
	local *gogh.LocalController,
	parser gogh.SpecParser,
	s string,
) func() error {
	return func() error {
		spec, alias, err := parser.ParseWithAlias(s)
		if err != nil {
			return err
		}

		token := tokens.Get(spec.Host(), spec.Owner())
		// check forked
		adaptor, err := github.NewAdaptor(ctx, spec.Host(), string(token))
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)
		repo, err := remote.Get(ctx, spec.Owner(), spec.Name(), nil)
		if err != nil {
			return err
		}

		l := log.FromContext(ctx).WithFields(log.Fields{
			"spec": spec,
		})
		l.Info("cloning")
		if _, err = local.Clone(ctx, spec, string(token), &gogh.LocalCloneOption{Alias: alias}); err != nil {
			l.WithField("error", err).Error("failed to get repository")
			return nil
		}
		if repo.Parent != nil && repo.Parent.String() != spec.String() {
			if err := local.SetRemoteSpecs(ctx, spec, map[string][]gogh.Spec{
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

func init() {
	cloneCommand.Flags().
		BoolVarP(&cloneFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(cloneCommand)
}
