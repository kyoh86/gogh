package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/gogh"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/infra/filesystem"
	"github.com/kyoh86/gogh/v4/infra/git"
	"github.com/kyoh86/gogh/v4/infra/github"
	"github.com/kyoh86/gogh/v4/infra/logger"
	"github.com/kyoh86/gogh/v4/ui/cli"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	ctx := logger.NewLogger(context.Background(), os.Stdout, os.Stderr)
	if err := run(ctx); err != nil {
		errString := err.Error()
		if strings.Contains(errString, "context canceled") || strings.Contains(errString, "context deadline exceeded") {
			// Ignore context cancellation errors
		} else {
			// Log the error with the Capitalized first letter
			log.FromContext(ctx).Errorf("%s: %s", gogh.AppName, strings.ToUpper(errString[:1])+errString[1:])
		}
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
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

	overlayStore := config.NewOverlayStore()
	overlayService, err := overlayStore.Load(ctx, func() overlay.OverlayService {
		return overlay.NewOverlayService(config.NewOverlayContentStore())
	})
	if err != nil {
		return fmt.Errorf("loading overlays: %w", err)
	}

	scriptStore := config.NewScriptStore()
	scriptService, err := scriptStore.Load(ctx, func() script.ScriptService {
		return script.NewScriptService(config.NewScriptSourceStore())
	})
	if err != nil {
		return fmt.Errorf("loading scripts: %w", err)
	}

	hookStore := config.NewHookStore()
	hookService, err := hookStore.Load(ctx, hook.NewHookService)
	if err != nil {
		return fmt.Errorf("loading hooks: %w", err)
	}

	svc := &service.ServiceSet{
		DefaultNameStore:   defaultNameStore,
		DefaultNameService: defaultNameService,

		TokenStore:   tokenStore,
		TokenService: tokenService,

		WorkspaceStore:   workspaceStore,
		WorkspaceService: workspaceService,

		OverlayStore:   overlayStore,
		OverlayService: overlayService,

		ScriptStore:   scriptStore,
		ScriptService: scriptService,

		HookStore:   hookStore,
		HookService: hookService,

		FlagsStore: flagsStore,
		Flags:      flags,

		ReferenceParser:     repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner()),
		HostingService:      github.NewHostingService(tokenService, defaultNameService),
		FinderService:       filesystem.NewFinderService(),
		AuthenticateService: github.NewAuthenticateService(),
		GitService:          git.NewService(git.CloneProgressWriter(os.Stdout)),
	}
	cmd, err := cli.NewApp(ctx, gogh.AppName, fmt.Sprintf("%s-%s (%s)", version, commit, date), svc)
	if err != nil {
		return fmt.Errorf("creating app: %w", err)
	}

	return cmd.ExecuteContext(ctx)
}
