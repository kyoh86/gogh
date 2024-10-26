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
	Short: "Login for the host and owner",
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
					Help: `You should generate personal access tokens with "Repository permissions":

- ✅ Read-only access to "Contents" and "Metadata"
- ✅ Read and write access to "Administration"`,
				},
			},
		}, &loginFlags); err != nil {
			return err
		}
		tokens.Set(loginFlags.Host, loginFlags.User, loginFlags.Password)
		return nil
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
	loginCommand.Flags().
		StringVarP(&loginFlags.Host, "host", "", gogh.DefaultHost, "Host name to login")
	loginCommand.Flags().StringVarP(&loginFlags.User, "user", "", "", "User name to login")
	loginCommand.Flags().
		StringVarP(&loginFlags.Password, "password", "", "", `Password or developer private token

You should generate personal access tokens with "Repository permissions":

- ✅ Read-only access to "Contents" and "Metadata"
- ✅ Read and write access to "Administration"`)
	authCommand.AddCommand(loginCommand)
}
