package commands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewRootsCommand(workspaceService workspace.WorkspaceService) *cobra.Command {
	return &cobra.Command{
		Use:     "roots",
		Short:   "Manage roots",
		Aliases: []string{"root"},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return config.SaveConfig()
		},
		Run: RootsListRun(workspaceService),
	}
}

func NewRootsListCommand(workspaceService workspace.WorkspaceService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all of the roots",
		Args:  cobra.ExactArgs(0),
		Run:   RootsListRun(workspaceService),
	}
}

func RootsListRun(workspaceService workspace.WorkspaceService) func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {
		for _, root := range workspaceService.GetRoots() {
			fmt.Println(root)
		}
	}
}

func NewRootsAddCommand(workspaceService workspace.WorkspaceService) *cobra.Command {
	var asPrimary bool
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add directories into the roots",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			return workspaceService.AddRoot(rootList[0], asPrimary)
		},
	}
	cmd.Flags().BoolVarP(&asPrimary, "as-primary", "", false, "Set as primary root")
	return cmd
}

func NewRootsRemoveCommand(workspaceService workspace.WorkspaceService) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a directory from the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				opts := make([]huh.Option[string], 0, len(workspaceService.GetRoots()))
				for _, root := range workspaceService.GetRoots() {
					opts = append(opts, huh.Option[string]{Key: root, Value: root})
				}

				form := huh.NewForm(huh.NewGroup(
					huh.NewSelect[string]().
						Title("Roots to remove").
						Options(opts...).
						Value(&selected),
				))
				if err := form.Run(); err != nil {
					return err
				}
			} else {
				selected = rootList[0]
			}
			return workspaceService.RemoveRoot(selected)
		},
	}
}

func NewRootsSetPrimaryCommand(workspaceService workspace.WorkspaceService) *cobra.Command {
	return &cobra.Command{
		Use:     "set-primary",
		Aliases: []string{"set-default"},
		Short:   "Set a directory as the primary in the roots",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				opts := make([]huh.Option[string], 0, len(workspaceService.GetRoots()))
				for _, root := range workspaceService.GetRoots() {
					opts = append(opts, huh.Option[string]{Key: root, Value: root})
				}

				form := huh.NewForm(huh.NewGroup(
					huh.NewSelect[string]().
						Title("A directory to set as default root").
						Options(opts...).
						Value(&selected),
				))
				if err := form.Run(); err != nil {
					return err
				}
			} else {
				selected = rootList[0]
			}

			return workspaceService.SetPrimaryRoot(selected)
		},
	}
}
