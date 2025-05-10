package cli

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/commands"
	"github.com/spf13/cobra"
)

func NewApp(
	ctx context.Context,
	conf *config.ConfigStore,
	defaultNameService repository.DefaultNameService,
	hostingService hosting.HostingService,
	finderService workspace.FinderService,
	workspaceService workspace.WorkspaceService,
	tokenService auth.TokenService,
	defaults *config.FlagStore,
) *cobra.Command {
	facadeCommand := &cobra.Command{
		Use:   config.AppName,
		Short: "GO GitHub local repository manager",
	}

	bundleCommand := commands.NewBundleCommand()
	bundleCommand.AddCommand(
		commands.NewBundleDumpCommand(conf, defaults, workspaceService, finderService),
		commands.NewBundleRestoreCommand(conf, defaultNameService, tokenService, defaults, hostingService, workspaceService),
	)

	authCommand := commands.NewAuthCommand(tokenService)
	authCommand.AddCommand(
		commands.NewAuthListCommand(tokenService),
		commands.NewAuthLoginCommand(tokenService),
		commands.NewAuthLogoutCommand(tokenService),
	)

	rootsCommand := commands.NewRootsCommand(conf)
	rootsCommand.AddCommand(
		commands.NewRootsSetPrimaryCommand(conf),
		commands.NewRootsRemoveCommand(conf),
		commands.NewRootsAddCommand(conf),
		commands.NewRootsListCommand(conf),
	)

	configCommand := commands.NewConfigCommand(conf, tokenService, defaults)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
	)

	facadeCommand.AddCommand(
		commands.NewManCommand(),
		commands.NewCwdCommand(conf, defaults, workspaceService, finderService),
		commands.NewListCommand(conf, defaults, workspaceService, finderService),
		commands.NewCloneCommand(conf, defaultNameService, tokenService, hostingService, workspaceService),
		commands.NewCreateCommand(conf, defaultNameService, tokenService, hostingService, workspaceService, defaults),
		commands.NewReposCommand(tokenService, hostingService, defaults),
		commands.NewDeleteCommand(conf, defaultNameService, tokenService, hostingService, finderService, workspaceService),
		commands.NewForkCommand(conf, defaultNameService, tokenService, defaults, hostingService),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)
	return facadeCommand
}
