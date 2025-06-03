package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/core/fs"
)

func TestReplaceTildeWithHome(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		path, err := testtarget.ReplaceTildeWithHome("")
		if err != nil {
			t.Fatalf("ReplaceTildeWithHome failed with error: %v", err)
		}
		if path != "" {
			t.Errorf("expected empty path, got %q", path)
		}
	})

	t.Run("tilde only", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("failed to get user home dir: %v", err)
		}

		path, err := testtarget.ReplaceTildeWithHome("~")
		if err != nil {
			t.Fatalf("ReplaceTildeWithHome failed with error: %v", err)
		}
		if path != homeDir {
			t.Errorf("expected %q, got %q", homeDir, path)
		}
	})

	t.Run("tilde with path", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("failed to get user home dir: %v", err)
		}
		expected := filepath.Join(homeDir, "test")

		path, err := testtarget.ReplaceTildeWithHome("~/test")
		if err != nil {
			t.Fatalf("ReplaceTildeWithHome failed with error: %v", err)
		}
		if path != expected {
			t.Errorf("expected %q, got %q", expected, path)
		}
	})

	t.Run("normal path", func(t *testing.T) {
		expected := "/normal/path"
		path, err := testtarget.ReplaceTildeWithHome("/normal/path")
		if err != nil {
			t.Fatalf("ReplaceTildeWithHome failed with error: %v", err)
		}
		if path != expected {
			t.Errorf("expected %q, got %q", expected, path)
		}
	})
}
