package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
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
	Args:  cobra.ExactArgs(0),
	RunE: func(*cobra.Command, []string) error {
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
		return servers.Set(loginFlags.Host, loginFlags.User, loginFlags.Password)
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

func init() {
	setup()
	loginCommand.Flags().StringVarP(&loginFlags.Host, "host", "", gogh.DefaultHost, "Host name to login")
	loginCommand.Flags().StringVarP(&loginFlags.User, "user", "", "", "User name to login")
	loginCommand.Flags().StringVarP(&loginFlags.Password, "password", "", "", "Password or developer private token")
	serversCommand.AddCommand(loginCommand)
}
