package cli

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(ctx context.Context) (*cobra.Command, error) {
	flags, flagsSource, err := store.LoadAlternative(
		ctx,
		config.DefaultFlags,
		config.NewFlagsStore(),
		config.NewFlagsStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load flags: %w", err)
	}

	defaultNameStore := config.NewDefaultNameStore()
	defaultNameService, defaultNameSource, err := store.LoadAlternative(
		ctx,
		config.DefaultName,
		defaultNameStore,
		config.NewDefaultNameStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default names: %w", err)
	}
	tokenStore := config.NewTokenStore()
	tokenService, tokenSource, err := store.LoadAlternative(
		ctx,
		config.DefaultTokenService,
		tokenStore,
		config.NewTokenStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	workspaceStore := config.NewWorkspaceStore()
	workspaceService, workspaceSource, err := store.LoadAlternative(
		ctx,
		config.DefaultWorkspaceService,
		workspaceStore,
		config.NewWorkspaceStoreV0(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load workspace: %w", err)
	}

	svc := commands.NewServiceSet(
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
		commands.NewBundleDumpCommand(svc),
		commands.NewBundleRestoreCommand(svc),
	)

	authCommand := commands.NewAuthCommand(svc)
	authCommand.AddCommand(
		commands.NewAuthListCommand(svc),
		commands.NewAuthLoginCommand(svc),
		commands.NewAuthLogoutCommand(svc),
	)

	rootsCommand := commands.NewRootsCommand(svc)
	rootsCommand.AddCommand(
		commands.NewRootsSetPrimaryCommand(svc),
		commands.NewRootsRemoveCommand(svc),
		commands.NewRootsAddCommand(svc),
		commands.NewRootsListCommand(svc),
	)

	configCommand := commands.NewConfigCommand(svc)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
		commands.NewSetDefaultHostCommand(svc),
		commands.NewSetDefaultOwnerCommand(svc),
	)

	appCommand.AddCommand(
		commands.NewMigrateCommand(svc, defaultNameStore, tokenStore, workspaceStore),
		commands.NewManCommand(),
		commands.NewCwdCommand(svc),
		commands.NewListCommand(svc),
		commands.NewCloneCommand(svc),
		commands.NewCreateCommand(svc),
		commands.NewReposCommand(svc),
		commands.NewDeleteCommand(svc),
		commands.NewForkCommand(svc),
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
