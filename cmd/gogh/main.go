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

	flags, flagsSource, err := store.LoadAlternative(
		ctx,
		config.DefaultFlags,
		config.NewFlagsStore(),
		config.NewFlagsStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("failed to load flags: %w", err)
	}

	defaultNameStore := config.NewDefaultNameStore()
	defaultNameService, defaultNameSource, err := store.LoadAlternative(
		ctx,
		config.DefaultName,
		defaultNameStore,
		config.NewDefaultNameStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("failed to load default names: %w", err)
	}
	tokenStore := config.NewTokenStore()
	tokenService, tokenSource, err := store.LoadAlternative(
		ctx,
		config.DefaultTokenService,
		tokenStore,
		config.NewTokenStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("failed to load tokens: %w", err)
	}

	workspaceStore := config.NewWorkspaceStore()
	workspaceService, workspaceSource, err := store.LoadAlternative(
		ctx,
		config.DefaultWorkspaceService,
		workspaceStore,
		config.NewWorkspaceStoreV0(),
	)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	cmd := cli.NewApp(
		defaultNameSource,
		defaultNameService,
		tokenSource,
		tokenService,
		workspaceSource,
		workspaceService,
		flagsSource,
		flags,
	)
	cmd.Version = fmt.Sprintf("%s-%s (%s)", version, commit, date)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return err
	}

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
