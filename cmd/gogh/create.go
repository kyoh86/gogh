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

var createFlags struct {
	template string
	dryrun   bool
}

var createCommand = &cobra.Command{
	Use:     "create [flags] [[OWNER/]NAME]",
	Aliases: []string{"new"},
	Short:   "Create a new project with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		servers := app.Servers()
		var name string
		if len(specs) == 0 {
			parser := gogh.NewSpecParser(servers)
			if err := survey.AskOne(&survey.Input{
				Message: "A spec of repository name to create",
			}, &name, survey.WithValidator(func(input interface{}) error {
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
			name = specs[0]
		}

		ctx := cmd.Context()
		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(name)
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

		if createFlags.template == "" {
			var ropt *gogh.RemoteCreateOption
			if server.User() != spec.Owner() {
				ropt = &gogh.RemoteCreateOption{Organization: spec.Owner()}
			}
			_, err = remote.Create(ctx, spec.Name(), ropt)
			return err
		}

		from, err := gogh.ParseSiblingSpec(spec, createFlags.template)
		if err != nil {
			return err
		}
		var ropt *gogh.RemoteCreateFromTemplateOption
		if server.User() != spec.Owner() {
			ropt = &gogh.RemoteCreateFromTemplateOption{Owner: spec.Owner()}
		}
		_, err = remote.CreateFromTemplate(ctx, from.Owner(), from.Name(), spec.Name(), ropt)
		return err
	},
}

func init() {
	createCommand.Flags().BoolVarP(&createFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	createCommand.Flags().StringVarP(&createFlags.template, "template", "", "", "Create new repository from the template")
	facadeCommand.AddCommand(createCommand)
}
