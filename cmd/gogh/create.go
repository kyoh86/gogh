package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/command"
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		servers := app.Servers()
		var selected string
		if len(specs) == 0 {
			configured, err := servers.List()
			if err != nil {
				return err
			}
			if len(configured) == 0 {
				return nil
			}
			specs = make([]string, 0, len(configured))
			for _, c := range configured {
				specs = append(specs, c.Host())
			}

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
		return command.Create(cmd.Context(), app.DefaultRoot(), app.Servers(), selected, nil, nil)
	},
}

func init() {
	facadeCommand.AddCommand(createCommand)
}
