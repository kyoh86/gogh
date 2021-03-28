package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var cloneCommand = &cobra.Command{
	Use:     "clone",
	Aliases: []string{"get"},
	Short:   "Clone a repository to local",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		ctx := cmd.Context()
		servers := app.Servers()
		var selected string
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
			if err := survey.AskOne(&survey.Select{
				Message: "A repository to clone",
				Options: specs,
			}, &selected); err != nil {
				return err
			}
		} else {
			selected = specs[0]
		}
		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(selected)
		if err != nil {
			return err
		}

		local := gogh.NewLocalController(app.DefaultRoot())
		_, err = local.Clone(ctx, spec, server, nil)
		return err
	},
}

func init() {
	facadeCommand.AddCommand(cloneCommand)
}
