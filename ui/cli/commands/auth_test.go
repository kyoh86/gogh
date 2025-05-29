package commands_test

import (
	"context"
	"testing"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/commands"
)

func TestNewAuthCommand(t *testing.T) {
	// Setup
	ctx := context.Background()
	serviceSet := &service.ServiceSet{}

	// Execute
	_, err := commands.NewAuthCommand(ctx, serviceSet)

	// Verify no error occurs
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
