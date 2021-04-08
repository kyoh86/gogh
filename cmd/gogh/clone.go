package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
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
	Use:     "clone",
	Aliases: []string{"get"},
	Short:   "Clone a repository to local",
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
		parser := gogh.NewSpecParser(servers)
		if cloneFlags.dryrun {
			for _, s := range specs {
				spec, _, err := parser.Parse(s)
				if err != nil {
					return err
				}

				p := gogh.NewProject(app.DefaultRoot(), spec)
				log.FromContext(ctx).Infof("git clone %q", p.URL())
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
	},
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

		l := log.FromContext(ctx).WithFields(log.Fields{
			"server": server,
			"spec":   spec,
		})
		l.Info("cloning")
		if _, err = local.Clone(ctx, spec, server, &gogh.LocalCloneOption{Alias: alias}); err != nil {
			l.WithField("error", err).Warn("failed to get repository")
			return nil
		}
		l.Info("finished")
		return nil
	}
}

func init() {
	cloneCommand.Flags().BoolVarP(&cloneFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(cloneCommand)
}
