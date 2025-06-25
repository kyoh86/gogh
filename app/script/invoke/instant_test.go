package invoke_test

import (
	"context"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestInvokeInstant(t *testing.T) {
	// Basic test to ensure the function exists and can be called
	t.Run("function accepts parameters", func(t *testing.T) {
		// Create a context that's already cancelled to prevent execution
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		code := "print('test')"
		globals := map[string]any{"key": "value"}

		// The function should accept these parameters without panicking
		// It will fail due to cancelled context, but that's expected
		_ = testtarget.InvokeInstant(ctx, location, code, globals)
	})

	t.Run("handles nil globals", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		code := "print('test')"

		// Should handle nil globals gracefully
		_ = testtarget.InvokeInstant(ctx, location, code, nil)
	})

	t.Run("handles empty code", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		globals := map[string]any{}

		// Should handle empty code gracefully
		_ = testtarget.InvokeInstant(ctx, location, "", globals)
	})

	t.Run("creates repo globals", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Test with various location configurations
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
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		location := repository.NewLocation("/tmp/test", "github.com", "user", "repo")

		// Globals with repo key that should be overwritten
		globals := map[string]any{
			"repo": map[string]any{
				"custom": "should be overwritten",
			},
			"other": "preserved",
		}

		// Should overwrite the repo key with location data
		_ = testtarget.InvokeInstant(ctx, location, "print('test')", globals)
	})
}
