package commands

import (
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

func NewAuthLoginCommand(tokens auth.TokenService) *cobra.Command {
	var f struct {
		Host string
	}

	useCase := auth_login.NewUseCase(tokens)

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"signin", "add"},
		Short:   "Login for the host and owner",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			deviceCodeCh, errCh := useCase.Execute(cmd.Context(), f.Host)
			for {
				select {
				case deviceCodeResp, ok := <-deviceCodeCh:
					if ok {
						if errors.Is(browser.OpenURL(deviceCodeResp.VerificationURI), exec.ErrNotFound) {
							fmt.Printf("Visit %s and enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
						} else {
							fmt.Printf("Opened %s, so enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
						}
					}
				case err, ok := <-errCh:
					if ok {
						return err
					}
					fmt.Println("Login successful!")
					return nil
				}
			}
		},
	}
	cmd.Flags().StringVarP(&f.Host, "host", "", "", "Host name to login")
	return cmd
}
