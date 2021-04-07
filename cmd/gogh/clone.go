package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
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
		var local *gogh.LocalController
		if !cloneFlags.dryrun {
			local = gogh.NewLocalController(app.DefaultRoot())
		}
		parser := gogh.NewSpecParser(servers)
		for _, s := range specs {
			spec, server, err := parser.Parse(s)
			if err != nil {
				return err
			}

			if cloneFlags.dryrun {
				p := gogh.NewProject(app.DefaultRoot(), spec)
				fmt.Printf("git clone %q\n", p.URL())
			} else {
				if _, err = local.Clone(ctx, spec, server, nil); err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	cloneCommand.Flags().BoolVarP(&cloneFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(cloneCommand)
}
