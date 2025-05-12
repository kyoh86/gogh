package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/infra/config"
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
	ctx := logger.NewLogger()

	flagsPathV0, err := config.FlagsPathV0()
	if err != nil {
		return fmt.Errorf("failed to get flags path (v0): %w", err)
	}
	flagsPath, err := config.FlagsPath()
	if err != nil {
		return fmt.Errorf("failed to get flags path: %w", err)
	}
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

	flags, err := store.LoadAlternative(
		ctx,
		config.DefaultFlags,
		config.NewFlagsStore(flagsPath),
		config.NewFlagsStoreV0(flagsPathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load flags: %w", err)
	}

	defaultNameStore := config.NewDefaultNameStore(defaultNamesPath)
	defaultNameService, err := store.LoadAlternative(
		ctx,
		config.DefaultName,
		defaultNameStore,
		config.NewDefaultNameStoreV0(tokensPathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load default names: %w", err)
	}
	tokenStore := config.NewTokenStore(tokensPath)
	tokenService, err := store.LoadAlternative(
		ctx,
		config.DefaultTokenService,
		tokenStore,
		config.NewTokenStoreV0(tokensPathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load tokens: %w", err)
	}

	workspaceStore := config.NewWorkspaceStore(workspacePath)
	workspaceService, err := store.LoadAlternative(
		ctx,
		config.DefaultWorkspaceService,
		workspaceStore,
		config.NewWorkspaceStoreV0(workspacePathV0),
	)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	cmd := cli.NewApp(
		defaultNameService,
		tokenService,
		workspaceService,
		flags,
	)
	cmd.Version = fmt.Sprintf("%s-%s (%s)", version, commit, date)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}

	if err := defaultNameStore.Save(ctx, defaultNameService); err != nil {
		return err
	}
	if err := tokenStore.Save(ctx, tokenService); err != nil {
		return err
	}
	if err := workspaceStore.Save(ctx, workspaceService); err != nil {
		return err
	}
	return nil
}
