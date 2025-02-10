package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/cli/browser"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var loginFlags struct {
	Host string
}

const clientID = "Ov23li6aEWIxek6F8P5L"

type TokenResponse struct {
	AccessToken string
	Scope       string
	TokenType   string
}

type ErrorResponse struct {
	Error            string
	ErrorDescription string
	ErrorURI         string
}

var loginCommand = &cobra.Command{
	Use:     "login",
	Aliases: []string{"signin", "add"},
	Short:   "Login for the host and owner",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginFlags.Host == "" {
			loginFlags.Host = github.DefaultHost
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().
					Title("Host name").
					Validate(gogh.ValidateHost).
					Value(&loginFlags.Host),
			)).Run(); err != nil {
				return err
			}
		}

		oauthConfig := &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:       fmt.Sprintf("https://%s/login/oauth/authorize", loginFlags.Host),
				TokenURL:      fmt.Sprintf("https://%s/login/oauth/access_token", loginFlags.Host),
				DeviceAuthURL: fmt.Sprintf("https://%s/login/device/code", loginFlags.Host),
			},
			Scopes: []string{"repo", "delete_repo"},
		}

		// Request device code
		deviceCodeResp, err := oauthConfig.DeviceAuth(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to request device code: %w", err)
		}

		if errors.Is(browser.OpenURL(deviceCodeResp.VerificationURI), exec.ErrNotFound) {
			fmt.Printf("Visit %s and enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
		} else {
			fmt.Printf("Opened %s, so enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
		}

		// Poll for token
		deviceCodeResp.Interval = deviceCodeResp.Interval + 1 // Add a second for safety
		tokenResp, err := oauthConfig.DeviceAccessToken(cmd.Context(), deviceCodeResp)
		if err != nil {
			return fmt.Errorf("failed to poll for token: %w", err)
		}

		// Get user info
		adaptor, err := github.NewAdaptor(context.Background(), loginFlags.Host, tokenResp.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to create GitHub adaptor: %w", err)
		}
		user, _, err := adaptor.GetAuthenticatedUser(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get authenticated user info: %w", err)
		}

		tokens.Set(loginFlags.Host, user.GetLogin(), tokenResp.AccessToken)

		fmt.Println("Login successful!")
		return nil
	},
}

func init() {
	loginCommand.Flags().StringVarP(&loginFlags.Host, "host", "", "", "Host name to login")
	authCommand.AddCommand(loginCommand)
}
