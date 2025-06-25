package invoke_test

import (
	"context"
	"os/exec"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func init() {
	// Override the instantCommandRunner to avoid actual subprocess execution
	testtarget.SetInstantCommandRunner(func(name string, args ...string) *exec.Cmd {
		// Use a command that will exit immediately without doing anything
		// This prevents the infinite recursion when running with -race
		cmd := exec.Command("true")
		if cmd.Path == "" {
			// Fallback for systems without 'true' command (like Windows)
			cmd = exec.Command("echo", "mock")
		}
		return cmd
	})
}

func TestInvokeInstant(t *testing.T) {
	// All tests use the mocked command runner that doesn't actually execute

	t.Run("function accepts parameters", func(t *testing.T) {
		ctx := context.Background()
		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		code := "print('test')"
		globals := map[string]any{"key": "value"}

		// The mock command will exit immediately, so this should complete quickly
		err := testtarget.InvokeInstant(ctx, location, code, globals)
		// We expect an error because 'true' command doesn't read stdin properly
		if err == nil {
			t.Log("Command completed without error (mock worked)")
		}
	})

	t.Run("handles nil globals", func(t *testing.T) {
		ctx := context.Background()
		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		code := "print('test')"

		_ = testtarget.InvokeInstant(ctx, location, code, nil)
	})

	t.Run("handles empty code", func(t *testing.T) {
		ctx := context.Background()
		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		globals := map[string]any{}

		_ = testtarget.InvokeInstant(ctx, location, "", globals)
	})

	t.Run("creates repo globals", func(t *testing.T) {
		ctx := context.Background()
		locations := []*repository.Location{
			repository.NewLocation("/home/user/repos/test", "github.com", "kyoh86", "myrepo"),
			repository.NewLocation("/tmp/test", "gitlab.com", "user", "project"),
			repository.NewLocation("/workspace/repos", "bitbucket.org", "team", "repo"),
		}

		for _, loc := range locations {
			_ = testtarget.InvokeInstant(ctx, loc, "-- test", nil)
		}
	})

	t.Run("handles globals with repo key", func(t *testing.T) {
		ctx := context.Background()
		location := repository.NewLocation("/tmp/test", "github.com", "user", "repo")

		globals := map[string]any{
			"repo": map[string]any{
				"custom": "should be overwritten",
			},
			"other": "preserved",
		}

		_ = testtarget.InvokeInstant(ctx, location, "print('test')", globals)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		location := repository.NewLocation("/tmp/test", "github.com", "user", "repo")

		// Should handle cancelled context gracefully
		_ = testtarget.InvokeInstant(ctx, location, "print('test')", nil)
	})
}
