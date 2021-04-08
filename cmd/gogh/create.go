package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{
	Use:     "create [flags] [[OWNER/]NAME]",
	Aliases: []string{"new"},
	Short:   "Create a new project with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		servers := app.Servers()
		var selected string
		if len(specs) == 0 {
			parser := gogh.NewSpecParser(servers)
			if err := survey.AskOne(&survey.Input{
				Message: "A spec of repository name to create",
			}, &selected, survey.WithValidator(func(input interface{}) error {
				s, ok := input.(string)
				if !ok {
					return errors.New("invalid type")
				}
				_, _, err := parser.Parse(s)
				return err
			})); err != nil {
				return err
			}
		} else {
			selected = specs[0]
		}

		ctx := cmd.Context()
		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(selected)
		if err != nil {
			return err
		}

		local := gogh.NewLocalController(app.DefaultRoot())
		if _, err = local.Create(ctx, spec, nil); err != nil {
			if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
				return err
			}
		}

		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)

		// check repo has already existed
		if _, err := remote.Get(ctx, spec.Owner(), spec.Name(), nil); err == nil {
			return nil
		}

		var ropt *gogh.RemoteCreateOption
		if server.User() != spec.Owner() {
			ropt = &gogh.RemoteCreateOption{Organization: spec.Owner()}
		}
		_, err = remote.Create(ctx, spec.Name(), ropt)
		return err
	},
}

func init() {
	facadeCommand.AddCommand(createCommand)
}
