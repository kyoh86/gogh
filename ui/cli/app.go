package cli

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/typ"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func newCmdWithSubs(
	ctx context.Context,
	svc *service.ServiceSet,
	main func(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error),
	fixedSubs []*cobra.Command,
	subs ...func(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error),
) (*cobra.Command, error) {
	cmd, err := main(ctx, svc)
	if err != nil {
		return nil, err
	}
	for _, sub := range fixedSubs {
		cmd.AddCommand(sub)
	}
	for _, subFn := range subs {
		subCmd, err := subFn(ctx, svc)
		if err != nil {
			return nil, err
		}
		cmd.AddCommand(subCmd)
	}
	return cmd, nil
}

func NewApp(
	ctx context.Context,
	appName string,
	version string,
	svc *service.ServiceSet,
) (*cobra.Command, error) {
	appCommand := &cobra.Command{
		Use:          appName,
		Short:        "GO GitHub local repository manager",
		SilenceUsage: true, // Do not show usage when error occurs; it is handled manually.
		PostRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if err := svc.DefaultNameStore.Save(ctx, svc.DefaultNameService, false); err != nil {
				return fmt.Errorf("saving default names: %w", err)
			}
			if err := svc.TokenStore.Save(ctx, svc.TokenService, false); err != nil {
				return fmt.Errorf("saving tokens: %w", err)
			}
			if err := svc.WorkspaceStore.Save(ctx, svc.WorkspaceService, false); err != nil {
				return fmt.Errorf("saving workspaces: %w", err)
			}
			return nil
		},
	}

	const (
		groupShow       = "show"
		groupManipulate = "manipulate"
		groupConfig     = "config"
	)
	appCommand.AddGroup(
		&cobra.Group{
			ID:    groupShow,
			Title: "Show repositories",
		},
		&cobra.Group{
			ID:    groupManipulate,
			Title: "Manipulate repositories",
		},
		&cobra.Group{
			ID:    groupConfig,
			Title: "Configurations",
		},
	)

	bundleCommand, err := newCmdWithSubs(
		ctx, svc,
		commands.NewBundleCommand,
		nil,
		commands.NewBundleDumpCommand,
		commands.NewBundleRestoreCommand,
	)
	if err != nil {
		return nil, err
	}

	authCommand, err := newCmdWithSubs(
		ctx, svc,
		commands.NewAuthCommand,
		nil,
		commands.NewAuthListCommand,
		commands.NewAuthLoginCommand,
		commands.NewAuthLogoutCommand,
	)
	if err != nil {
		return nil, err
	}
	authCommand.GroupID = groupConfig

	rootsCommand, err := newCmdWithSubs(
		ctx, svc,
		commands.NewRootsCommand,
		nil,
		commands.NewRootsSetPrimaryCommand,
		commands.NewRootsRemoveCommand,
		commands.NewRootsAddCommand,
		commands.NewRootsListCommand,
	)
	if err != nil {
		return nil, err
	}
	rootsCommand.GroupID = groupConfig

	configAuthCommand := typ.Ptr(*authCommand)
	configAuthCommand.GroupID = ""
	configRootsCommand := typ.Ptr(*rootsCommand)
	configRootsCommand.GroupID = ""
	configCommand, err := newCmdWithSubs(
		ctx, svc,
		commands.NewConfigCommand,
		[]*cobra.Command{configAuthCommand, configRootsCommand},
		commands.NewSetDefaultHostCommand,
		commands.NewSetDefaultOwnerCommand,
	)
	if err != nil {
		return nil, err
	}
	configCommand.GroupID = groupConfig

	cmds := []*cobra.Command{
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	}
	for _, sub := range []struct {
		fn    func(context.Context, *service.ServiceSet) (*cobra.Command, error)
		group string
	}{
		{fn: commands.NewMigrateCommand, group: groupConfig},
		{fn: commands.NewManCommand},
		{fn: commands.NewCwdCommand, group: groupShow},
		{fn: commands.NewListCommand, group: groupShow},
		{fn: commands.NewCloneCommand, group: groupManipulate},
		{fn: commands.NewCreateCommand, group: groupManipulate},
		{fn: commands.NewReposCommand, group: groupShow},
		{fn: commands.NewDeleteCommand, group: groupManipulate},
		{fn: commands.NewForkCommand, group: groupManipulate},
	} {
		c, err := sub.fn(ctx, svc)
		if err != nil {
			return nil, err
		}
		c.GroupID = sub.group
		cmds = append(cmds, c)
	}
	appCommand.AddCommand(cmds...)

	return appCommand, nil
}
