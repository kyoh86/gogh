package commands

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/cli/browser"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
)

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

func NewAuthLoginCommand(tokens *config.TokenManager) *cobra.Command {
	var f struct {
		Host string
	}

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"signin", "add"},
		Short:   "Login for the host and owner",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if f.Host == "" {
				f.Host = github.DefaultHost
				if err := huh.NewForm(huh.NewGroup(
					huh.NewInput().
						Title("Host name").
						Validate(gogh.ValidateHost).
						Value(&f.Host),
				)).Run(); err != nil {
					return err
				}
			}

			oauthConfig := github.OAuth2Config(f.Host)

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
			deviceCodeResp.Interval++ // Add a second for safety
			tokenResp, err := oauthConfig.DeviceAccessToken(cmd.Context(), deviceCodeResp)
			if err != nil {
				return fmt.Errorf("failed to poll for token: %w", err)
			}

			if tokenResp == nil {
				return fmt.Errorf("got nil token response")
			}

			// Get user info
			adaptor, err := github.NewAdaptor(context.Background(), f.Host, tokenResp)
			if err != nil {
				return fmt.Errorf("failed to create GitHub adaptor: %w", err)
			}
			user, _, err := adaptor.GetAuthenticatedUser(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get authenticated user info: %w", err)
			}

			tokens.Set(f.Host, user.GetLogin(), *tokenResp)

			fmt.Println("Login successful!")
			return nil
		},
	}
	cmd.Flags().StringVarP(&f.Host, "host", "", "", "Host name to login")
	return cmd
}
