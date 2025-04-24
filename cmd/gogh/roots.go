package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/config"
	"github.com/spf13/cobra"
)

func NewRootsCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "roots",
		Short:   "Manage roots",
		Aliases: []string{"root"},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return config.SaveConfig()
		},
		Run: RootsListRun(conf),
	}
}

func NewRootsListCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all of the roots",
		Args:  cobra.ExactArgs(0),
		Run:   RootsListRun(conf),
	}
}

func RootsListRun(conf *config.Config) func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {
		for _, root := range conf.GetRoots() {
			fmt.Println(root)
		}
	}
}

func NewRootsAddCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add directories into the roots",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			return conf.AddRoots(rootList)
		},
	}
}

func NewRootsRemoveCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a directory from the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				options := make([]huh.Option[string], 0, len(conf.GetRoots()))
				for _, root := range conf.GetRoots() {
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
			conf.RemoveRoot(selected)
			return nil
		},
	}
}

func NewRootsSetDefaultCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "set-default",
		Short: "Set a directory as the default in the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				options := make([]huh.Option[string], 0, len(conf.GetRoots()))
				for _, root := range conf.GetRoots() {
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

			return conf.SetDefaultRoot(selected)
		},
	}
}
