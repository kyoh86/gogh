package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var authListCommand = &cobra.Command{
	Use:   "list",
	Short: "Listup authenticated host and owners",
	Args:  cobra.ExactArgs(0),
	RunE: func(*cobra.Command, []string) error {
		for _, entry := range tokens.Entries() {
			fmt.Println(entry)
		}
		return nil
	},
}

func init() {
	authCommand.AddCommand(authListCommand)
}
