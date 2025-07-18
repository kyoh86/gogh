package cli

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/commands"
	"github.com/spf13/cobra"
)

func cmdWithSubs(
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
		Version:      version,
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
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
			if err := svc.FlagsStore.Save(ctx, svc.Flags, false); err != nil {
				return fmt.Errorf("saving flags: %w", err)
			}
			if err := svc.OverlayStore.Save(ctx, svc.OverlayService, false); err != nil {
				return fmt.Errorf("saving overlays: %w", err)
			}
			if err := svc.ScriptStore.Save(ctx, svc.ScriptService, false); err != nil {
				return fmt.Errorf("saving scripts: %w", err)
			}
			if err := svc.HookStore.Save(ctx, svc.HookService, false); err != nil {
				return fmt.Errorf("saving hooks: %w", err)
			}
			if err := svc.ExtraStore.Save(ctx, svc.ExtraService, false); err != nil {
				return fmt.Errorf("saving extra: %w", err)
			}
			return nil
		},
	}

	const (
		groupShow       = "show"
		groupManipulate = "manipulate"
		groupAutomation = "automation"
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
			ID:    groupAutomation,
			Title: "Automation",
		},
		&cobra.Group{
			ID:    groupConfig,
			Title: "Configurations",
		},
	)

	bundleCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewBundleCommand,
		nil,
		commands.NewBundleDumpCommand,
		commands.NewBundleRestoreCommand,
	)
	if err != nil {
		return nil, err
	}

	authCommand, err := cmdWithSubs(
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

	rootsCommand, err := cmdWithSubs(
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

	overlayCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewOverlayCommand,
		nil,
		commands.NewOverlayAddCommand,
		commands.NewOverlayApplyCommand,
		commands.NewOverlayEditCommand,
		commands.NewOverlayListCommand,
		commands.NewOverlayRemoveCommand,
		commands.NewOverlayShowCommand,
		commands.NewOverlayUpdateCommand,
	)
	if err != nil {
		return nil, err
	}
	overlayCommand.GroupID = groupAutomation

	scriptCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewScriptCommand,
		nil,
		commands.NewScriptAddCommand,
		commands.NewScriptCreateCommand,
		commands.NewScriptInvokeCommand,
		commands.NewScriptEditCommand,
		commands.NewScriptListCommand,
		commands.NewScriptRemoveCommand,
		commands.NewScriptShowCommand,
		commands.NewScriptRunCommand,
		commands.NewScriptInvokeInstantCommand,
		commands.NewScriptUpdateCommand,
	)
	if err != nil {
		return nil, err
	}
	scriptCommand.GroupID = groupAutomation

	hookCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewHookCommand,
		nil,
		commands.NewHookAddCommand,
		commands.NewHookInvokeCommand,
		commands.NewHookListCommand,
		commands.NewHookRemoveCommand,
		commands.NewHookShowCommand,
		commands.NewHookUpdateCommand,
	)
	if err != nil {
		return nil, err
	}
	hookCommand.GroupID = groupAutomation

	extraCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewExtraCommand,
		nil,
		commands.NewExtraSaveCommand,
		commands.NewExtraCreateCommand,
		commands.NewExtraListCommand,
		commands.NewExtraShowCommand,
		commands.NewExtraRemoveCommand,
		commands.NewExtraApplyCommand,
	)
	if err != nil {
		return nil, err
	}
	extraCommand.GroupID = groupAutomation

	configAuthCommand := typ.Ptr(*authCommand)
	configAuthCommand.GroupID = ""
	configRootsCommand := typ.Ptr(*rootsCommand)
	configRootsCommand.GroupID = ""
	configCommand, err := cmdWithSubs(
		ctx, svc,
		commands.NewConfigCommand,
		[]*cobra.Command{configAuthCommand, configRootsCommand},
		commands.NewConfigShowCommand,
		commands.NewSetDefaultHostCommand,
		commands.NewSetDefaultOwnerCommand,
		commands.NewMigrateCommand,
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
		overlayCommand,
		scriptCommand,
		hookCommand,
		extraCommand,
	}
	for _, sub := range []struct {
		fn    func(context.Context, *service.ServiceSet) (*cobra.Command, error)
		group string
	}{
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
