package commands

import (
	"errors"
	"fmt"

	"github.com/kyoh86/gogh/v3/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
)

func NewAuthSetDefaultCommand(tokens *config.TokenManager) *cobra.Command {
	var f struct {
		Host  string
		Owner string
	}

	cmd := &cobra.Command{
		Use:   "set-default",
		Short: "Set a host and an owner as the default in the auth",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			if f.Host == "" {
				return errors.New("host is required")
			}
			if f.Owner == "" {
				tokens.DefaultHost = f.Host
				return nil
			}
			if err := tokens.SetDefaultOwner(f.Host, f.Owner); err != nil {
				return fmt.Errorf("failed to set the default owner: %w", err)
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&f.Host, "host", "", github.DefaultHost, "Host name to login")
	flags.StringVarP(&f.Owner, "owner", "", "", "Owner name to login")
	return cmd
}
