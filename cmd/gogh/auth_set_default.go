package main

import (
	"errors"
	"fmt"

	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var setDefaultFlags struct {
	Host  string
	Owner string
}

var authSetDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "Set a host and an owner as the default in the auth",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, rootList []string) error {
		if setDefaultFlags.Host == "" {
			return errors.New("host is required")
		}
		if setDefaultFlags.Owner == "" {
			tokens.DefaultHost = setDefaultFlags.Host
			return nil
		}
		if err := tokens.SetDefaultOwner(setDefaultFlags.Host, setDefaultFlags.Owner); err != nil {
			return fmt.Errorf("failed to set the default owner: %w", err)
		}
		return nil
	},
}

func init() {
	flags := authSetDefaultCommand.Flags()
	flags.AddFlagSet(&pflag.FlagSet{})
	flags.StringVarP(&setDefaultFlags.Host, "host", "", github.DefaultHost, "Host name to login")
	flags.StringVarP(&setDefaultFlags.Owner, "owner", "", "", "Owner name to login")
	authCommand.AddCommand(authSetDefaultCommand)
}
