package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var rootsCommand = &cobra.Command{
	Use:     "roots",
	Short:   "Manage roots",
	Aliases: []string{"root"},
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return app.SaveConfig()
	},
	Run: rootsListCommand.Run,
}

var rootsListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all of the roots",
	Args:  cobra.ExactArgs(0),
	Run: func(*cobra.Command, []string) {
		for _, root := range app.Roots() {
			fmt.Println(root)
		}
	},
}

var rootsAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add directories into the roots",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, roots []string) {
		app.AddRoots(roots)
	},
}

var rootsRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a directory from the roots",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, roots []string) error {
		var selected string
		if len(roots) == 0 {
			if err := survey.AskOne(&survey.Select{
				Message: "Roots to remove",
				Options: app.Roots(),
			}, &selected); err != nil {
				return err
			}
		} else {
			selected = roots[0]
		}
		app.RemoveRoot(selected)
		return nil
	},
}

var rootsSetDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "Set a directory as the default in the roots",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, roots []string) error {
		var selected string
		if len(roots) == 0 {
			if err := survey.AskOne(&survey.Select{
				Message: "A directory to set as default root",
				Options: app.Roots(),
			}, &selected); err != nil {
				return err
			}
		} else {
			selected = roots[0]
		}
		app.SetDefaultRoot(selected)
		return nil
	},
}

func init() {
	rootsCommand.AddCommand(rootsSetDefaultCommand)
	rootsCommand.AddCommand(rootsRemoveCommand)
	rootsCommand.AddCommand(rootsAddCommand)
	rootsCommand.AddCommand(rootsListCommand)
	facadeCommand.AddCommand(rootsCommand)
}
