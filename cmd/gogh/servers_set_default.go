package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var setDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "Set default server",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, hosts []string) error {
		var selected string
		if len(hosts) == 0 {
			configured, err := servers.List()
			if err != nil {
				return err
			}
			if len(configured) == 0 {
				return nil
			}
			hosts = make([]string, 0, len(configured))
			for _, c := range configured {
				hosts = append(hosts, c.Host())
			}

			if err := survey.AskOne(&survey.Select{
				Message: "A server to set as default",
				Options: hosts,
			}, &selected); err != nil {
				return err
			}
		} else {
			selected = hosts[0]
		}
		return servers.SetDefault(selected)
	},
}

func init() {
	setup()
	serversCommand.AddCommand(setDefaultCommand)
}
