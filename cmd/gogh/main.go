package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/infra/logger"
	"github.com/kyoh86/gogh/v3/ui/cli"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func loadConfigOrExit[T any](name string, loader func() (T, error)) T {
	v, err := loader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load %s: %s\n", name, err)
		os.Exit(1)
	}
	return v
}

func main() {
	if err := run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx := logger.NewLogger()

	conf := loadConfigOrExit("config", config.LoadConfig)
	defaults := loadConfigOrExit("flags", config.LoadFlags)

	tokensPathV0, err := config.TokensPathV0()
	if err != nil {
		return fmt.Errorf("failed to get tokens path (v0): %w", err)
	}
	tokensPath, err := config.TokensPath()
	if err != nil {
		return fmt.Errorf("failed to get tokens path (v0): %w", err)
	}
	workspacePath, err := config.WorkspacePath()
	if err != nil {
		return fmt.Errorf("failed to get workspace path: %w", err)
	}
	workspacePathV0, err := config.WorkspacePathV0()
	if err != nil {
		return fmt.Errorf("failed to get workspace path (v0): %w", err)
	}
	defaultNamesPath, err := config.DefaultNamesPath()
	if err != nil {
		return fmt.Errorf("failed to get default names path: %w", err)
	}

	defaultNameService, err := store.LoadAlternative(ctx,
		config.NewDefaultNameStore(defaultNamesPath),
		config.NewDefaultNameStoreV0(tokensPathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load default names: %w", err)
	}

	tokenService, err := store.LoadAlternative(ctx,
		config.NewTokenStore(tokensPath),
		config.NewTokenStoreV0(tokensPathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load tokens: %w", err)
	}

	workspaceService, err := store.LoadAlternative(ctx,
		config.NewWorkspaceStore(workspacePath),
		config.NewWorkspaceStoreV0(workspacePathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	hostingService := github.NewHostingService(tokenService)
	finderService := filesystem.NewFinderService()

	cmd := cli.NewApp(ctx, conf, defaultNameService, hostingService, finderService, workspaceService, tokenService, defaults)
	cmd.Version = fmt.Sprintf("%s-%s (%s)", version, commit, date)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil
}
