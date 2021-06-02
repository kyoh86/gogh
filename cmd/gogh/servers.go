package main

import (
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

var (
	servers     gogh.Servers
	serversPath string
)

func Servers() *gogh.Servers {
	return &servers
}

func setupServers() error {
	serversPath = filepath.Join(cacheDir, Name, "servers.yaml")
	if err := loadYAML(serversPath, &servers); err != nil {
		return fmt.Errorf("load servers: %w", err)
	}
	return nil
}

func SaveServers() error {
	if err := saveYAML(serversPath, &servers); err != nil {
		return fmt.Errorf("save servers: %w", err)
	}
	return nil
}

var serversCommand = &cobra.Command{
	Use:     "servers",
	Short:   "Manage servers",
	Aliases: []string{"server"},
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return SaveServers()
	},
}

var setDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "Set default server",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, hosts []string) error {
		servers := Servers()
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
	serversCommand.AddCommand(setDefaultCommand)
	facadeCommand.AddCommand(serversCommand)
}
