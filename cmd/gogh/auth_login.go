package main

import (
	"errors"

	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:    "Ov23li6aEWIxek6F8P5L",
		RedirectURL: "http://localhost",
		Endpoint:    github.Endpoint,
		Scopes:      []string{"repo"},
	}
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		authCodeURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

		fmt.Printf("Visit the URL for the auth dialog: %v\n", authCodeURL)

		var authCode string
		if err := survey.AskOne(&survey.Input{
			Message: "Enter the authorization code:",
		}, &authCode); err != nil {
			return err
		}

		token, err := oauthConfig.Exchange(ctx, authCode)
		if err != nil {
			return fmt.Errorf("failed to exchange token: %w", err)
		}

		tokens.Set(loginFlags.Host, loginFlags.User, token.AccessToken)
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
