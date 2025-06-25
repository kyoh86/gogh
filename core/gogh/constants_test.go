package gogh_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/core/gogh"
)

func TestConstants(t *testing.T) {
	// Test DefaultHost constant
	if testtarget.DefaultHost != "github.com" {
		t.Errorf("DefaultHost = %q, want %q", testtarget.DefaultHost, "github.com")
	}

	// Test AppName constant
	if testtarget.AppName != "gogh" {
		t.Errorf("AppName = %q, want %q", testtarget.AppName, "gogh")
	}
}

func TestConstantsUsage(t *testing.T) {
	// Demonstrate usage of constants
	t.Logf("Application name: %s", testtarget.AppName)
	t.Logf("Default host: %s", testtarget.DefaultHost)

	// Verify constants are exported and accessible
	_ = testtarget.DefaultHost
	_ = testtarget.AppName
}
