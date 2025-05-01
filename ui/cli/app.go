package cli

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(ctx context.Context, conf *config.ConfigStore, defaultNameService repository.DefaultNameService, tokens auth.TokenService, defaults *config.FlagStore) *cobra.Command {
	facadeCommand := &cobra.Command{
		Use:   config.AppName,
		Short: "GO GitHub local repository manager",
	}

	bundleCommand := commands.NewBundleCommand()
	bundleCommand.AddCommand(
		commands.NewBundleDumpCommand(conf, defaults),
		commands.NewBundleRestoreCommand(conf, defaultNameService, tokens, defaults),
	)

	authCommand := commands.NewAuthCommand(tokens)
	authCommand.AddCommand(
		commands.NewAuthListCommand(tokens),
		commands.NewAuthLoginCommand(tokens),
		commands.NewAuthLogoutCommand(tokens),
	)

	rootsCommand := commands.NewRootsCommand(conf)
	rootsCommand.AddCommand(
		commands.NewRootsSetDefaultCommand(conf),
		commands.NewRootsRemoveCommand(conf),
		commands.NewRootsAddCommand(conf),
		commands.NewRootsListCommand(conf),
	)

	configCommand := commands.NewConfigCommand(conf, tokens, defaults)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
	)

	facadeCommand.AddCommand(
		commands.NewCwdCommand(conf, defaults),
		commands.NewListCommand(conf, defaults),
		commands.NewCloneCommand(conf, defaultNameService, tokens),
		commands.NewCreateCommand(conf, defaultNameService, tokens, defaults),
		commands.NewReposCommand(tokens, defaults),
		commands.NewDeleteCommand(conf, defaultNameService, tokens),
		commands.NewForkCommand(conf, defaultNameService, tokens, defaults),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)
	return facadeCommand
}
