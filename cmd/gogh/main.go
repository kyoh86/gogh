package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/gogh"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	"github.com/kyoh86/gogh/v3/infra/git"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/infra/logger"
	"github.com/kyoh86/gogh/v3/ui/cli"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	if err := run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx := logger.NewLogger(context.Background())

	flagsStore := config.NewFlagsStore()
	flags, err := config.LoadAlternative(
		ctx,
		config.DefaultFlags,
		flagsStore,
		config.NewFlagsStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("loading flags: %w", err)
	}

	defaultNameStore := config.NewDefaultNameStore()
	defaultNameService, err := config.LoadAlternative(
		ctx,
		repository.NewDefaultNameService,
		defaultNameStore,
		config.NewDefaultNameStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("loading default names: %w", err)
	}

	tokenStore := config.NewTokenStore()
	tokenService, err := config.LoadAlternative(
		ctx,
		auth.NewTokenService,
		tokenStore,
		config.NewTokenStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("loading tokens: %w", err)
	}

	workspaceStore := config.NewWorkspaceStore()
	workspaceService, err := config.LoadAlternative(
		ctx,
		filesystem.NewWorkspaceService,
		workspaceStore,
		config.NewWorkspaceStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("loading workspace: %w", err)
	}

	svc := &service.ServiceSet{
		DefaultNameStore:   defaultNameStore,
		DefaultNameService: defaultNameService,

		TokenStore:   tokenStore,
		TokenService: tokenService,

		WorkspaceStore:   workspaceStore,
		WorkspaceService: workspaceService,

		FlagsStore: flagsStore,
		Flags:      flags,

		ReferenceParser:     repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner()),
		HostingService:      github.NewHostingService(tokenService),
		FinderService:       filesystem.NewFinderService(),
		AuthenticateService: github.NewAuthenticateService(),
		GitService:          git.NewService(),
	}
	cmd, err := cli.NewApp(ctx, gogh.AppName, fmt.Sprintf("%s-%s (%s)", version, commit, date), svc)
	if err != nil {
		return fmt.Errorf("creating app: %w", err)
	}

	return cmd.ExecuteContext(ctx)
}
