package cli

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(
	ctx context.Context,
	svc *service.ServiceSet,
) (*cobra.Command, error) {
	appCommand := &cobra.Command{
		Use:   config.AppName,
		Short: "GO GitHub local repository manager",
	}

	bundleCommand := commands.NewBundleCommand()
	bundleCommand.AddCommand(
		commands.NewBundleDumpCommand(ctx, svc),
		commands.NewBundleRestoreCommand(ctx, svc),
	)

	authCommand := commands.NewAuthCommand(ctx, svc)
	authCommand.AddCommand(
		commands.NewAuthListCommand(ctx, svc),
		commands.NewAuthLoginCommand(ctx, svc),
		commands.NewAuthLogoutCommand(ctx, svc),
	)

	rootsCommand := commands.NewRootsCommand(ctx, svc)
	rootsCommand.AddCommand(
		commands.NewRootsSetPrimaryCommand(ctx, svc),
		commands.NewRootsRemoveCommand(ctx, svc),
		commands.NewRootsAddCommand(ctx, svc),
		commands.NewRootsListCommand(ctx, svc),
	)

	configCommand := commands.NewConfigCommand(ctx, svc)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
		commands.NewSetDefaultHostCommand(ctx, svc),
		commands.NewSetDefaultOwnerCommand(ctx, svc),
	)

	appCommand.AddCommand(
		commands.NewMigrateCommand(ctx, svc),
		commands.NewManCommand(),
		commands.NewCwdCommand(ctx, svc),
		commands.NewListCommand(ctx, svc),
		commands.NewCloneCommand(ctx, svc),
		commands.NewCreateCommand(ctx, svc),
		commands.NewReposCommand(ctx, svc),
		commands.NewDeleteCommand(ctx, svc),
		commands.NewForkCommand(ctx, svc),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)

	appCommand.PostRunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := svc.DefaultNameStore.Save(ctx, svc.DefaultNameService, false); err != nil {
			return err
		}
		if err := svc.TokenStore.Save(ctx, svc.TokenService, false); err != nil {
			return err
		}
		if err := svc.WorkspaceStore.Save(ctx, svc.WorkspaceService, false); err != nil {
			return err
		}
		return nil
	}
	return appCommand, nil
}
