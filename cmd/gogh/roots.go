package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var rootsCommand = &cobra.Command{
	Use:     "roots",
	Short:   "Manage roots",
	Aliases: []string{"root"},
	PersistentPostRunE: func(*cobra.Command, []string) error {
		return saveConfig()
	},
	Run: rootsListCommand.Run,
}

var rootsListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all of the roots",
	Args:  cobra.ExactArgs(0),
	Run: func(*cobra.Command, []string) {
		for _, root := range roots() {
			fmt.Println(root)
		}
	},
}

var rootsAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add directories into the roots",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, rootList []string) error {
		return addRoots(rootList)
	},
}

var rootsRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a directory from the roots",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, rootList []string) error {
		var selected string
		if len(rootList) == 0 {
			options := make([]huh.Option[string], 0, len(roots()))
			for _, root := range roots() {
				options = append(options, huh.Option[string]{Key: root, Value: root})
			}

			form := huh.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("Roots to remove").
					Options(options...).
					Value(&selected),
			))
			if err := form.Run(); err != nil {
				return err
			}
		} else {
			selected = rootList[0]
		}
		removeRoot(selected)
		return nil
	},
}

var rootsSetDefaultCommand = &cobra.Command{
	Use:   "set-default",
	Short: "Set a directory as the default in the roots",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(_ *cobra.Command, rootList []string) error {
		var selected string
		if len(rootList) == 0 {
			options := make([]huh.Option[string], 0, len(roots()))
			for _, root := range roots() {
				options = append(options, huh.Option[string]{Key: root, Value: root})
			}

			form := huh.NewForm(huh.NewGroup(
				huh.NewSelect[string]().
					Title("A directory to set as default root").
					Options(options...).
					Value(&selected),
			))
			if err := form.Run(); err != nil {
				return err
			}
		} else {
			selected = rootList[0]
		}

		return setDefaultRoot(selected)
	},
}

func init() {
	rootsCommand.AddCommand(rootsSetDefaultCommand)
	rootsCommand.AddCommand(rootsRemoveCommand)
	rootsCommand.AddCommand(rootsAddCommand)
	rootsCommand.AddCommand(rootsListCommand)
	configCommand.AddCommand(rootsCommand)
	facadeCommand.AddCommand(rootsCommand)
}
