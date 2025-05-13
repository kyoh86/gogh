package commands

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func NewRootsCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:     "roots",
		Short:   "Manage roots",
		Aliases: []string{"root"},
		Run:     RootsListRun(svc),
	}
}

func NewRootsListCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all of the roots",
		Args:  cobra.NoArgs,
		Run:   RootsListRun(svc),
	}
}

func RootsListRun(svc *ServiceSet) func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {
		for _, root := range svc.workspaceService.GetRoots() {
			fmt.Println(root)
		}
	}
}

func NewRootsAddCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	var asPrimary bool
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add directories into the roots",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			return svc.workspaceService.AddRoot(rootList[0], asPrimary)
		},
	}
	cmd.Flags().BoolVarP(&asPrimary, "as-primary", "", false, "Set as primary root")
	return cmd
}

func NewRootsRemoveCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a directory from the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				opts := make([]huh.Option[string], 0, len(svc.workspaceService.GetRoots()))
				for _, root := range svc.workspaceService.GetRoots() {
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
			return svc.workspaceService.RemoveRoot(selected)
		},
	}
}

func NewRootsSetPrimaryCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:     "set-primary",
		Aliases: []string{"set-default"},
		Short:   "Set a directory as the primary in the roots",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			var selected string
			if len(rootList) == 0 {
				opts := make([]huh.Option[string], 0, len(svc.workspaceService.GetRoots()))
				for _, root := range svc.workspaceService.GetRoots() {
					opts = append(opts, huh.Option[string]{Key: root, Value: root})
				}

				form := huh.NewForm(huh.NewGroup(
					huh.NewSelect[string]().
						Title("A directory to set as primary root").
						Options(opts...).
						Value(&selected),
				))
				if err := form.Run(); err != nil {
					return err
				}
			} else {
				selected = rootList[0]
			}

			return svc.workspaceService.SetPrimaryRoot(selected)
		},
	}
}
