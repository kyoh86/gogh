package cli

import (
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(conf *config.Config, tokens *config.TokenManager, defaults *config.Flags) *cobra.Command {
	facadeCommand := &cobra.Command{
		Use:   config.AppName,
		Short: "GO GitHub project manager",
	}

	bundleCommand := commands.NewBundleCommand()
	bundleCommand.AddCommand(
		commands.NewBundleDumpCommand(conf, defaults),
		commands.NewBundleRestoreCommand(conf, tokens, defaults),
	)

	authCommand := commands.NewAuthCommand()
	authCommand.AddCommand(
		commands.NewAuthListCommand(tokens),
		commands.NewAuthLoginCommand(tokens),
		commands.NewAuthLogoutCommand(tokens),
		commands.NewAuthSetDefaultCommand(tokens),
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
		commands.NewCloneCommand(conf, tokens),
		commands.NewCreateCommand(conf, tokens, defaults),
		commands.NewReposCommand(tokens, defaults),
		commands.NewDeleteCommand(conf, tokens),
		commands.NewForkCommand(conf, tokens, defaults),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)
	return facadeCommand
}
