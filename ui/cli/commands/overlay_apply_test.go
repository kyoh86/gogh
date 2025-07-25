package commands_test

import (
	"context"
	"testing"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/commands"
)

func TestNewOverlayApplyCommand(t *testing.T) {
	// Setup
	ctx := context.Background()
	serviceSet := &service.ServiceSet{Flags: &config.Flags{}}

	// Execute and verify no error occurs
	_, err := commands.NewOverlayApplyCommand(ctx, serviceSet)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
