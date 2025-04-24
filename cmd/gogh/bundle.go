package main

import (
	"github.com/spf13/cobra"
)

func NewBundleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "bundle",
		Short: "Manage bundle",
	}
}
