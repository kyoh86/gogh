package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewRootsCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "roots",
		Short:   "Manage roots",
		Aliases: []string{"root"},
		RunE:    RootsListRunE(svc),
	}, nil
}

func NewRootsListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "list",
		Short: "List all of the roots",
		Args:  cobra.NoArgs,
		RunE:  RootsListRunE(svc),
	}, nil
}

func RootsListRunE(svc *service.ServiceSet) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		roots := svc.WorkspaceService.GetRoots()
		if len(roots) == 0 {
			return errors.New("no roots found: you need to set root by `gogh roots add`")
		}
		for _, root := range roots {
			fmt.Println(root)
		}
		return nil
	}
}

func NewRootsAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var asPrimary bool
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add directories into the roots",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, rootList []string) error {
			ctx := cmd.Context()
			if err := svc.WorkspaceService.AddRoot(rootList[0], asPrimary); err != nil {
				return err
			}
			log.FromContext(ctx).Infof("Added root: %q", rootList[0])
			return nil
		},
	}
	flags.BoolVarP(cmd, &asPrimary, "as-primary", "", false, "Set as primary root")
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

func NewRootsRemoveCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a directory from the roots",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, rootList []string) error {
			selected, err := selectRoot(svc, "Roots to remove", rootList)
			if err != nil {
				return err
			}
			if err := svc.WorkspaceService.RemoveRoot(selected); err != nil {
				return err
			}
			log.FromContext(ctx).Infof("Removed root: %q", selected)
			return nil
		},
	}, nil
}

func NewRootsSetPrimaryCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "set-primary",
		Aliases: []string{"set-default"},
		Short:   "Set a directory as the primary in the roots",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, rootList []string) error {
			selected, err := selectRoot(svc, "A directory to set as primary root", rootList)
			if err != nil {
				return err
			}
			if err := svc.WorkspaceService.SetPrimaryRoot(selected); err != nil {
				return err
			}
			log.FromContext(cmd.Context()).Infof("Set %q as primary root", selected)
			return nil
		},
	}, nil
}
