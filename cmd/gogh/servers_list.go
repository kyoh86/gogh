package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serversListCommand = &cobra.Command{
	Use:   "list",
	Short: "Listup servers",
	Args:  cobra.ExactArgs(0),
	RunE: func(*cobra.Command, []string) error {
		list, err := servers.List()
		if err != nil {
			return fmt.Errorf("listup servers: %w", err)
		}
		for _, server := range list {
			fmt.Println(server)
		}
		return nil
	},
}

func init() {
	setup()
	serversCommand.AddCommand(serversListCommand)
}
