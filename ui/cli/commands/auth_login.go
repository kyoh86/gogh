package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/cli/browser"
	"github.com/kyoh86/gogh/v3/app/auth_login"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewAuthLoginCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		Host string
	}

	useCase := auth_login.NewUseCase(svc.TokenService, svc.AuthenticateService, svc.HostingService)

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"signin", "add"},
		Short:   "Login for the host and owner",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if f.Host == "" {
				f.Host = svc.DefaultNameService.GetDefaultHost()
				if err := huh.NewForm(huh.NewGroup(
					huh.NewInput().
						Title("Host name").
						Value(&f.Host),
				)).Run(); err != nil {
					return err
				}
			}

			if err := useCase.Execute(cmd.Context(), f.Host, func(ctx context.Context, res auth_login.DeviceAuthResponse) error {
				if errors.Is(browser.OpenURL(res.VerificationURI), exec.ErrNotFound) {
					fmt.Fprintf(
						os.Stderr,
						"Failed to open browser automatically. Please open this URL in your browser:\n%s\n\nThen enter the code: %s\n",
						res.VerificationURI,
						res.UserCode,
					)
				} else {
					fmt.Fprintf(
						os.Stderr,
						"Your browser has been opened to: %s\nPlease enter this code in the browser: %s\n",
						res.VerificationURI,
						res.UserCode,
					)
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
	return cmd, nil
}
