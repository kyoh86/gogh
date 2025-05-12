package commands

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/cli/browser"
	"github.com/kyoh86/gogh/v3/app/auth_login"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/repository"
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

func NewAuthLoginCommand(svc *ServiceSet) *cobra.Command {
	var f struct {
		Host string
	}

	useCase := auth_login.NewUseCase(svc.tokenService, svc.authenticateService, svc.hostingService)

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"signin", "add"},
		Short:   "Login for the host and owner",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if f.Host == "" {
				f.Host = github.GlobalHost
				if err := huh.NewForm(huh.NewGroup(
					huh.NewInput().
						Title("Host name").
						Validate(repository.ValidateHost).
						Value(&f.Host),
				)).Run(); err != nil {
					return err
				}
			}

			if err := useCase.Execute(cmd.Context(), f.Host, func(ctx context.Context, response auth.DeviceAuthResponse) error {
				if errors.Is(browser.OpenURL(response.VerificationURI), exec.ErrNotFound) {
					fmt.Printf("Visit %s and enter the code: %s\n", response.VerificationURI, response.UserCode)
				} else {
					fmt.Printf("Opened %s, so enter the code: %s\n", response.VerificationURI, response.UserCode)
				}
				return nil
			}); err != nil {
				return err
			}
			fmt.Println("Login successful!")
			return nil
		},
	}
	cmd.Flags().StringVarP(&f.Host, "host", "", "", "Host name to login")
	return cmd
}
