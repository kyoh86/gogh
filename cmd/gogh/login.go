package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var loginFlags struct {
	Host     string
	User     string
	Password string
}

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Login for a server",
	RunE: func(*cobra.Command, []string) error {
		if err := getParams(); err != nil {
			return err
		}
		servers := app.Servers()
		if err := servers.Set(loginFlags.Host, loginFlags.User, loginFlags.Password); err != nil {
			return err
		}
		return app.SaveServers()
	},
}

func stringValidator(v func(s string) error) survey.Validator {
	return func(i interface{}) error {
		s, ok := i.(string)
		if !ok {
			return errors.New("invalid type")
		}
		return v(s)
	}
}

func getParams() error {
	var _ survey.Select
	if err := survey.Ask([]*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Host name",
				Default: loginFlags.Host,
			},
			Validate: stringValidator(gogh.ValidateHost),
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "User name",
				Default: loginFlags.User,
			},
			Validate: stringValidator(gogh.ValidateOwner),
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password or developer private token",
			},
		},
	}, &loginFlags); err != nil {
		return err
	}
	return nil
}

func init() {
	loginCommand.Flags().StringVarP(&loginFlags.Host, "host", "", gogh.DefaultHost, "Host name to login")
	loginCommand.Flags().StringVarP(&loginFlags.User, "user", "", "", "User name to login")
	loginCommand.Flags().StringVarP(&loginFlags.Password, "password", "", "", "Password or developer private token")
	facadeCommand.AddCommand(loginCommand)
}
