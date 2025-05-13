package cli

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(ctx context.Context) (*cobra.Command, error) {
	flags, flagsSource, err := config.LoadAlternative(
		ctx,
		config.DefaultFlags,
		config.NewFlagsStore(),
		config.NewFlagsStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load flags: %w", err)
	}

	defaultNameStore := config.NewDefaultNameStore()
	defaultNameService, defaultNameSource, err := config.LoadAlternative(
		ctx,
		config.DefaultName,
		defaultNameStore,
		config.NewDefaultNameStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default names: %w", err)
	}
	tokenStore := config.NewTokenStore()
	tokenService, tokenSource, err := config.LoadAlternative(
		ctx,
		config.DefaultTokenService,
		tokenStore,
		config.NewTokenStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	workspaceStore := config.NewWorkspaceStore()
	workspaceService, workspaceSource, err := config.LoadAlternative(
		ctx,
		config.DefaultWorkspaceService,
		workspaceStore,
		config.NewWorkspaceStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load workspace: %w", err)
	}

	svc := service.NewServiceSet(
		defaultNameSource,
		defaultNameService,
		tokenSource,
		tokenService,
		workspaceSource,
		workspaceService,
		flagsSource,
		flags,
	)
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
		commands.NewMigrateCommand(ctx, svc, defaultNameStore, tokenStore, workspaceStore),
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
		if err := defaultNameStore.Save(ctx, defaultNameService, false); err != nil {
			return err
		}
		if err := tokenStore.Save(ctx, tokenService, false); err != nil {
			return err
		}
		if err := workspaceStore.Save(ctx, workspaceService, false); err != nil {
			return err
		}
		return nil
	}
	return appCommand, nil
}
