package invoke

import (
	"bytes"
	"os/exec"
	"testing"
)

// Test the execCmd implementation
func TestExecCmd(t *testing.T) {
	// Test with a simple command that should exist on all platforms
	ec := &execCmd{exec.Command("echo", "test")}

	t.Run("SetDir", func(t *testing.T) {
		ec.SetDir("/tmp")
		if ec.Dir != "/tmp" {
			t.Errorf("SetDir failed: expected /tmp, got %s", ec.Dir)
		}
	})

	t.Run("SetStdout", func(t *testing.T) {
		var buf bytes.Buffer
		ec.SetStdout(&buf)
		if ec.Stdout != &buf {
			t.Error("SetStdout failed")
		}
	})

	t.Run("SetStderr", func(t *testing.T) {
		var buf bytes.Buffer
		ec.SetStderr(&buf)
		if ec.Stderr != &buf {
			t.Error("SetStderr failed")
		}
	})

	t.Run("StdinPipe", func(t *testing.T) {
		// Create a new command for this test
		cmd := &execCmd{exec.Command("echo", "test")}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("StdinPipe failed: %v", err)
		}
		if stdin == nil {
			t.Error("StdinPipe returned nil")
		}
		stdin.Close()
	})

	t.Run("Run", func(t *testing.T) {
		// Use a simple command that should work on all platforms
		cmd := &execCmd{exec.Command("echo", "hello")}
		var out bytes.Buffer
		cmd.SetStdout(&out)

		err := cmd.Run()
		if err != nil {
			t.Errorf("Run failed: %v", err)
		}

		// Check that something was written to stdout
		if out.Len() == 0 {
			t.Error("Expected output from echo command")
		}
	})
}

// Test defaultCommandRunner
func TestDefaultCommandRunner(t *testing.T) {
	cmd := defaultCommandRunner("echo", "test")

	// Verify it returns a Command interface
	if cmd == nil {
		t.Fatal("defaultCommandRunner returned nil")
	}

	// Verify it's actually an execCmd
	ec, ok := cmd.(*execCmd)
	if !ok {
		t.Fatal("defaultCommandRunner did not return an execCmd")
	}

	// Verify the command was set correctly
	if ec.Path != "echo" && !contains(ec.Path, "echo") {
		t.Errorf("Expected command path to contain 'echo', got %s", ec.Path)
	}

	if len(ec.Args) < 2 || ec.Args[1] != "test" {
		t.Errorf("Expected args to contain 'test', got %v", ec.Args)
	}
}

// Test SetCommandRunner
func TestSetCommandRunner_ReturnsPreviousRunner(t *testing.T) {
	// Save original
	original := commandRunner
	defer func() {
		commandRunner = original
	}()

	// Create a custom runner
	customCalled := false
	customRunner := func(name string, args ...string) Command {
		customCalled = true
		return &execCmd{exec.Command(name, args...)}
	}

	// Set the custom runner and verify it returns the previous one
	previous := SetCommandRunner(customRunner)

	// The previous should be the original defaultCommandRunner
	if previous == nil {
		t.Error("SetCommandRunner should return the previous runner")
	}

	// Verify our custom runner is now active
	cmd := commandRunner("test", "arg")
	if !customCalled {
		t.Error("Custom runner was not called")
	}
	if cmd == nil {
		t.Error("Custom runner returned nil")
	}

	// Restore using the returned function
	SetCommandRunner(previous)

	// Reset the flag
	customCalled = false

	// Verify it's restored
	_ = commandRunner("test", "arg")
	if customCalled {
		t.Error("Custom runner should not be called after restore")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
