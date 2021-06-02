package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serversListFlags struct{}

var serversListCommand = &cobra.Command{
	Use:   "list",
	Short: "Listup servers",
	Args:  cobra.ExactArgs(0),
	RunE: func(*cobra.Command, []string) error {
		servers, err := Servers().List()
		if err != nil {
			return fmt.Errorf("listup servers: %w", err)
		}
		for _, server := range servers {
			fmt.Println(server)
		}
		return nil
	},
}

func init() {
	serversCommand.AddCommand(serversListCommand)
}
