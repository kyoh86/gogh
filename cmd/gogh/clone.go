package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
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
  They'll be filled with a server-spec set by "servers login".

  It accepts an alias for each repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example"
  into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"

  Note that a host name in the alias is ignored:
    "github.com/kyoh86/example=example.com/kyoh86-tryouts/sample"
      will be placed in
    "$(gogh root)/github.com/kyoh86-tryouts/sample"`,

	RunE: func(cmd *cobra.Command, specs []string) error {
		ctx := cmd.Context()
		servers := app.Servers()
		if len(specs) == 0 {
			servers, err := servers.List()
			if err != nil {
				return err
			}
			for _, server := range servers {
				adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
				if err != nil {
					return err
				}
				remote := gogh.NewRemoteController(adaptor)
				founds, err := remote.List(ctx, nil)
				if err != nil {
					return err
				}
				for _, s := range founds {
					specs = append(specs, s.String())
				}
			}
			if err := survey.AskOne(&survey.MultiSelect{
				Message: "A repository to clone",
				Options: specs,
			}, &specs); err != nil {
				return err
			}
		}
		return cloneAll(ctx, servers, specs, cloneFlags.dryrun)
	},
}

func cloneAll(ctx context.Context, servers *gogh.Servers, specs []string, dryrun bool) error {
	parser := gogh.NewSpecParser(servers)
	if dryrun {
		for _, s := range specs {
			spec, _, err := parser.Parse(s)
			if err != nil {
				return err
			}

			log.FromContext(ctx).Infof("git clone %q", spec.URL())
		}
		return nil
	}

	local := gogh.NewLocalController(app.DefaultRoot())
	if len(specs) == 1 {
		return cloneOne(ctx, local, parser, specs[0])()
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, s := range specs {
		eg.Go(cloneOne(ctx, local, parser, s))
	}
	return eg.Wait()
}

func cloneOne(ctx context.Context, local *gogh.LocalController, parser *gogh.SpecParser, s string) func() error {
	return func() error {
		var alias *gogh.Spec
		part := strings.Split(s, "=")
		switch len(part) {
		case 1:
			// noop
		case 2:
			as, _, err := parser.Parse(part[1])
			if err != nil {
				return err
			}
			alias = &as
			s = part[0]
		default:
			return fmt.Errorf("invalid spec: %s", s)
		}
		spec, server, err := parser.Parse(s)
		if err != nil {
			return err
		}

		// check forked
		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)
		parent, err := remote.GetParent(ctx, spec.Owner(), spec.Name(), nil)
		if err != nil {
			return err
		}

		l := log.FromContext(ctx).WithFields(log.Fields{
			"server": server,
			"spec":   spec,
		})
		l.Info("cloning")
		if _, err = local.Clone(ctx, spec, server, &gogh.LocalCloneOption{Alias: alias}); err != nil {
			l.WithField("error", err).Warn("failed to get repository")
			return nil
		}
		if parent.String() != spec.String() {
			if err := local.SetRemoteSpecs(ctx, spec, map[string][]gogh.Spec{
				git.DefaultRemoteName: {spec},
				"upstream":            {parent},
			}); err != nil {
				return err
			}
		}
		l.Info("finished")
		return nil
	}
}

func init() {
	cloneCommand.Flags().BoolVarP(&cloneFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(cloneCommand)
}
