package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var logoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "Logout from a server",
	RunE: func(cmd *cobra.Command, hosts []string) error {
		servers := app.Servers()
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

			var selected []string
			if err := survey.AskOne(&survey.MultiSelect{
				Message: "Hosts to logout from",
				Options: hosts,
			}, &selected); err != nil {
				return err
			}
			hosts = selected
		}

		for _, host := range hosts {
			log.FromContext(cmd.Context()).WithField("host", host).Info("logout from")
			if err := servers.Remove(host); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	serversCommand.AddCommand(logoutCommand)
}
