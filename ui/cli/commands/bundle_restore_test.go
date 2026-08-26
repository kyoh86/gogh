package commands_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/commands"
)

func setHome(t *testing.T, home string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
		t.Setenv("HOME", "")
		return
	}
	t.Setenv("HOME", home)
}

func TestNewBundleRestoreCommand(t *testing.T) {
	// Setup
	ctx := context.Background()
	serviceSet := &service.ServiceSet{Flags: &config.Flags{}}

	// Execute and verify no error occurs
	_, err := commands.NewBundleRestoreCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestBundleRestoreCommandExpandsTildeFile(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	setHome(t, home)

	bundlePath := filepath.Join(home, "bundle.txt")
	if err := os.WriteFile(bundlePath, []byte("github.com/kyoh86/gogh\n"), 0o644); err != nil {
		t.Fatalf("write bundle file: %v", err)
	}

	serviceSet := &service.ServiceSet{
		Flags: &config.Flags{
			BundleRestore: config.BundleRestoreFlags{
				CloneRetryLimit: 1,
			},
		},
	}
	cmd, err := commands.NewBundleRestoreCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cmd.SetArgs([]string{"--dry-run", "--file", "~/bundle.txt"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
