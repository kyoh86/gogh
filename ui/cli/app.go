package cli

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/git"
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
	authService auth.AuthenticateService,
	defaults *config.FlagStore,
	gitService git.GitService,
) *cobra.Command {
	facadeCommand := &cobra.Command{
		Use:   config.AppName,
		Short: "GO GitHub local repository manager",
	}

	bundleCommand := commands.NewBundleCommand()
	bundleCommand.AddCommand(
		commands.NewBundleDumpCommand(defaults, workspaceService, finderService, gitService),
		commands.NewBundleRestoreCommand(defaultNameService, tokenService, defaults, hostingService, workspaceService, gitService),
	)

	authCommand := commands.NewAuthCommand(tokenService)
	authCommand.AddCommand(
		commands.NewAuthListCommand(tokenService),
		commands.NewAuthLoginCommand(tokenService, authService, hostingService),
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
		commands.NewCloneCommand(defaultNameService, tokenService, hostingService, workspaceService, gitService),
		commands.NewCreateCommand(defaultNameService, tokenService, hostingService, workspaceService, defaults, gitService),
		commands.NewReposCommand(tokenService, hostingService, defaults),
		commands.NewDeleteCommand(defaultNameService, tokenService, hostingService, finderService, workspaceService),
		commands.NewForkCommand(defaultNameService, tokenService, defaults, hostingService, workspaceService, gitService),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)
	return facadeCommand
}
