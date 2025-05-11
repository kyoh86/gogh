package cli

import (
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(svc *commands.ServiceSet) *cobra.Command {
	facadeCommand := &cobra.Command{
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
	)

	facadeCommand.AddCommand(
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
	return facadeCommand
}
