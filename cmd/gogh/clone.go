package main

import (
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
		} else {
			local := gogh.NewLocalController(app.DefaultRoot())
			eg, ctx := errgroup.WithContext(ctx)
			for _, s := range specs {
				s := s
				eg.Go(func() error {
					spec, server, err := parser.Parse(s)
					if err != nil {
						return err
					}

					log.FromContext(ctx).WithFields(log.Fields{
						"server": server,
						"spec":   spec,
					}).Info("cloning")
					if _, err = local.Clone(ctx, spec, server, nil); err != nil {
						log.FromContext(ctx).WithField("error", err).Warn("failed to get repository")
						return nil
					}
					log.FromContext(ctx).WithFields(log.Fields{
						"server": server,
						"spec":   spec,
					}).Info("finished")
					return nil
				})
			}
			return eg.Wait()
		}
		return nil
	},
}

func init() {
	cloneCommand.Flags().BoolVarP(&cloneFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(cloneCommand)
}
