package commands

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewRootsCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "roots",
		Short:   "Manage roots",
		Aliases: []string{"root"},
		Run:     RootsListRun(svc),
	}, nil
}

func NewRootsListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "list",
		Short: "List all of the roots",
		Args:  cobra.NoArgs,
		Run:   RootsListRun(svc),
	}, nil
}

func RootsListRun(svc *service.ServiceSet) func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {
		for _, root := range svc.WorkspaceService.GetRoots() {
			fmt.Println(root)
		}
	}
}

func NewRootsAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var asPrimary bool
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add directories into the roots",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			return svc.WorkspaceService.AddRoot(rootList[0], asPrimary)
		},
	}
	cmd.Flags().BoolVarP(&asPrimary, "as-primary", "", false, "Set as primary root")
	return cmd, nil
}

func selectRoot(svc *service.ServiceSet, title string, rootList []string) (string, error) {
	if len(rootList) > 0 {
		return rootList[0], nil
	}
	var selected string
	opts := make([]huh.Option[string], 0, len(svc.WorkspaceService.GetRoots()))
	for _, root := range svc.WorkspaceService.GetRoots() {
		opts = append(opts, huh.Option[string]{Key: root, Value: root})
	}
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title(title).
			Options(opts...).
			Value(&selected),
	))
	if err := form.Run(); err != nil {
		return "", err
	}
	return selected, nil
}

func NewRootsRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a directory from the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			selected, err := selectRoot(svc, "Roots to remove", rootList)
			if err != nil {
				return err
			}
			return svc.WorkspaceService.RemoveRoot(selected)
		},
	}, nil
}

func NewRootsSetPrimaryCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "set-primary",
		Aliases: []string{"set-default"},
		Short:   "Set a directory as the primary in the roots",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			selected, err := selectRoot(svc, "A directory to set as primary root", rootList)
			if err != nil {
				return err
			}
			return svc.WorkspaceService.SetPrimaryRoot(selected)
		},
	}, nil
}
