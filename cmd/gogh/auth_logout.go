package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
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
				indices = append(indices, fmt.Sprintf("%s/%s", c.Host, c.Owner))
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
			words := strings.SplitN(target, "/", 2)
			if len(words) != 2 {
				log.FromContext(cmd.Context()).WithField("target", target).Error("invalid target (must be host/owner)")
				continue
			}
			tokens.Delete(words[0], words[1])
		}
		return nil
	},
}

func init() {
	authCommand.AddCommand(logoutCommand)
}
