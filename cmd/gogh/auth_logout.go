package main

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/config"
	"github.com/spf13/cobra"
)

func NewAuthLogoutCommand(tokens *config.TokenManager) *cobra.Command {
	return &cobra.Command{
		Use:     "logout",
		Aliases: []string{"signout", "remove"},
		Short:   "Logout from the host and owner",
		RunE: func(cmd *cobra.Command, indices []string) error {
			if len(indices) == 0 {
				configured := tokens.Entries()
				if len(configured) == 0 {
					return nil
				}
				options := make([]huh.Option[string], 0, len(configured))
				for _, c := range configured {
					name := fmt.Sprintf("%s/%s", c.Host, c.Owner)
					options = append(options, huh.Option[string]{Key: name, Value: name})
				}

				var selected []string
				if err := huh.NewForm(huh.NewGroup(
					huh.NewMultiSelect[string]().
						Title("Hosts to logout from").
						Options(options...).
						Value(&selected),
				)).Run(); err != nil {
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
}
