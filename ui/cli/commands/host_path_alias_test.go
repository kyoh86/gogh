package commands_test

import (
	"context"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/infra/filesystem"
	"github.com/kyoh86/gogh/v4/ui/cli/commands"
)

func TestNewHostPathAliasCommand(t *testing.T) {
	ctx := context.Background()
	serviceSet := &service.ServiceSet{Flags: &config.Flags{}}

	_, err := commands.NewHostPathAliasCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHostPathAliasSetAndRemoveCommand(t *testing.T) {
	ctx := context.Background()
	workspaceService := filesystem.NewWorkspaceService()
	serviceSet := &service.ServiceSet{
		Flags:            &config.Flags{},
		WorkspaceService: workspaceService,
	}

	setCmd, err := commands.NewHostPathAliasSetCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	setCmd.SetContext(ctx)
	if err := setCmd.RunE(setCmd, []string{"github.com", "gh"}); err != nil {
		t.Fatalf("Expected set command to succeed, got %v", err)
	}
	aliases := workspaceService.GetHostPathAliases()
	if got := aliases["github.com"]; got != "gh" {
		t.Fatalf("Expected github.com alias to be gh, got %q", got)
	}

	removeCmd, err := commands.NewHostPathAliasRemoveCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	removeCmd.SetContext(ctx)
	if err := removeCmd.RunE(removeCmd, []string{"github.com"}); err != nil {
		t.Fatalf("Expected remove command to succeed, got %v", err)
	}
	aliases = workspaceService.GetHostPathAliases()
	if _, ok := aliases["github.com"]; ok {
		t.Fatalf("Expected github.com alias to be removed")
	}
}

func TestHostPathAliasSetCommandRejectsPathAlias(t *testing.T) {
	ctx := context.Background()
	serviceSet := &service.ServiceSet{
		Flags:            &config.Flags{},
		WorkspaceService: filesystem.NewWorkspaceService(),
	}

	cmd, err := commands.NewHostPathAliasSetCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if err := cmd.RunE(cmd, []string{"github.com", "github/short"}); err == nil {
		t.Fatalf("Expected set command to reject alias with path separator")
	}
}
