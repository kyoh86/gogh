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
		commands.NewBundleDumpCommand(defaults, workspaceService, finderService),
		commands.NewBundleRestoreCommand(defaultNameService, tokenService, defaults, hostingService, workspaceService),
	)

	authCommand := commands.NewAuthCommand(tokenService)
	authCommand.AddCommand(
		commands.NewAuthListCommand(tokenService),
		commands.NewAuthLoginCommand(tokenService),
		commands.NewAuthLogoutCommand(tokenService),
	)

	rootsCommand := commands.NewRootsCommand(workspaceService)
	rootsCommand.AddCommand(
		commands.NewRootsSetPrimaryCommand(workspaceService),
		commands.NewRootsRemoveCommand(workspaceService),
		commands.NewRootsAddCommand(workspaceService),
		commands.NewRootsListCommand(workspaceService),
	)

	configCommand := commands.NewConfigCommand(workspaceService, tokenService, defaults)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
	)

	facadeCommand.AddCommand(
		commands.NewManCommand(),
		commands.NewCwdCommand(defaults, workspaceService, finderService),
		commands.NewListCommand(defaults, workspaceService, finderService),
		commands.NewCloneCommand(defaultNameService, tokenService, hostingService, workspaceService),
		commands.NewCreateCommand(defaultNameService, tokenService, hostingService, workspaceService, defaults),
		commands.NewReposCommand(tokenService, hostingService, defaults),
		commands.NewDeleteCommand(defaultNameService, tokenService, hostingService, finderService, workspaceService),
		commands.NewForkCommand(defaultNameService, tokenService, defaults, hostingService),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)
	return facadeCommand
}
