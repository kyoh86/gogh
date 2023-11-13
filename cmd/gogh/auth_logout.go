package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

var logoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the host and owner",
	RunE: func(cmd *cobra.Command, indices []string) error {
		if len(indices) == 0 {
			configured := tokens.Entries()
			if len(configured) == 0 {
				return nil
			}
			indices = make([]string, 0, len(configured))
			for _, c := range configured {
				indices = append(indices, c.String())
			}

			var selected []string
			if err := survey.AskOne(&survey.MultiSelect{
				Message: "Hosts to logout from",
				Options: indices,
			}, &selected); err != nil {
				return err
			}
			indices = selected
		}

		for _, target := range indices {
			log.FromContext(cmd.Context()).WithField("target", target).Info("logout from")
			target, err := gogh.ParseTokenTarget(target)
			if err != nil {
				log.FromContext(cmd.Context()).WithField("target", target).Error("invalid target")
				continue
			}
			tokens.Delete(target.Host, target.Owner)
		}
		return nil
	},
}

func init() {
	authCommand.AddCommand(logoutCommand)
}
