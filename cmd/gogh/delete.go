package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove"},
	Short:   "Delete a repository with a remote repository",
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
				Message: "A repository to delete",
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
		if err := local.Delete(ctx, spec, nil); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("delete local: %w", err)
			}
		}

		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		return gogh.NewRemoteController(adaptor).Delete(ctx, spec.Owner(), spec.Name(), nil)
	},
}

func init() {
	facadeCommand.AddCommand(deleteCommand)
}
