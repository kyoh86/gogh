package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
)

func TestFlags_HasChanges(t *testing.T) {
	f := &testtarget.Flags{}
	if f.HasChanges() {
		t.Errorf("HasChanges should always return false")
	}
}

func TestFlags_MarkSaved(t *testing.T) {
	f := &testtarget.Flags{}
	// No assertions needed, just ensuring it doesn't panic
	f.MarkSaved()
}

func TestDefaultFlags(t *testing.T) {
	f := testtarget.DefaultFlags()
	if f == nil {
		t.Fatalf("DefaultFlags() returned nil")
	}

	homeDir, err := os.UserHomeDir()
	if err == nil && homeDir != "" {
		expectedPath := filepath.Join(homeDir, "./.config/gogh/bundle.txt")
		if f.BundleDump.File != expectedPath {
			t.Errorf("expected BundleDump.File to be %q, got %q", expectedPath, f.BundleDump.File)
		}
		if f.BundleRestore.File != expectedPath {
			t.Errorf("expected BundleRestore.File to be %q, got %q", expectedPath, f.BundleRestore.File)
		}
	}

	// Check default values
	if f.BundleRestore.CloneRetryLimit != 3 {
		t.Errorf("expected BundleRestore.CloneRetryLimit to be 3, got %d", f.BundleRestore.CloneRetryLimit)
	}
	if f.BundleRestore.RequestTimeout != 5*time.Minute {
		t.Errorf("expected BundleRestore.RequestTimeout to be 5m, got %v", f.BundleRestore.RequestTimeout)
	}
	if f.Clone.RequestTimeout != 5*time.Minute {
		t.Errorf("expected Clone.RequestTimeout to be 5m, got %v", f.Clone.RequestTimeout)
	}
	if f.Repos.Limit != 30 {
		t.Errorf("expected Repos.Limit to be 30, got %d", f.Repos.Limit)
	}
	if f.Repos.Color != "auto" {
		t.Errorf("expected Repos.Color to be 'auto', got %q", f.Repos.Color)
	}
	if !reflect.DeepEqual(f.Repos.Relation, []string{"owner", "organization-member"}) {
		t.Errorf("expected Repos.Relation to be ['owner', 'organization-member'], got %v", f.Repos.Relation)
	}
	if f.Create.RequestTimeout != 5*time.Minute {
		t.Errorf("expected Create.RequestTimeout to be 5m, got %v", f.Create.RequestTimeout)
	}
	if f.Create.CloneRetryLimit != 3 {
		t.Errorf("expected Create.CloneRetryLimit to be 3, got %d", f.Create.CloneRetryLimit)
	}
	if f.List.Limit != 100 {
		t.Errorf("expected List.Limit to be 100, got %d", f.List.Limit)
	}
	if f.Fork.RequestTimeout != 5*time.Minute {
		t.Errorf("expected Fork.RequestTimeout to be 5m, got %v", f.Fork.RequestTimeout)
	}
	if f.Fork.CloneRetryLimit != 3 {
		t.Errorf("expected Fork.CloneRetryLimit to be 3, got %d", f.Fork.CloneRetryLimit)
	}

	// Default boolean flags should be false
	if f.BundleRestore.Dryrun {
		t.Errorf("expected BundleRestore.Dryrun to be false")
	}
	if f.Clone.Dryrun {
		t.Errorf("expected Clone.Dryrun to be false")
	}
	if f.Create.Dryrun {
		t.Errorf("expected Create.Dryrun to be false")
	}
	if f.List.Primary {
		t.Errorf("expected List.Primary to be false")
	}
	if f.Fork.DefaultBranchOnly {
		t.Errorf("expected Fork.DefaultBranchOnly to be false")
	}
}

func TestFlagStructsInitialization(t *testing.T) {
	// Test that all flag structs can be initialized
	bundleDump := testtarget.BundleDumpFlags{}
	if bundleDump.File != "" {
		t.Errorf("expected empty File, got %q", bundleDump.File)
	}

	bundleRestore := testtarget.BundleRestoreFlags{}
	if bundleRestore.RequestTimeout != 0 {
		t.Errorf("expected zero RequestTimeout, got %v", bundleRestore.RequestTimeout)
	}
	if bundleRestore.File != "" {
		t.Errorf("expected empty File, got %q", bundleRestore.File)
	}
	if bundleRestore.CloneRetryLimit != 0 {
		t.Errorf("expected zero CloneRetryLimit, got %d", bundleRestore.CloneRetryLimit)
	}
	if bundleRestore.Dryrun {
		t.Errorf("expected false Dryrun")
	}

	clone := testtarget.CloneFlags{}
	if clone.RequestTimeout != 0 {
		t.Errorf("expected zero RequestTimeout, got %v", clone.RequestTimeout)
	}
	if clone.Dryrun {
		t.Errorf("expected false Dryrun")
	}

	create := testtarget.CreateFlags{}
	if create.Template != "" {
		t.Errorf("expected empty Template, got %q", create.Template)
	}
	if create.Private {
		t.Errorf("expected false Private")
	}

	cwd := testtarget.CwdFlags{}
	if cwd.Format != "" {
		t.Errorf("expected empty Format, got %q", cwd.Format)
	}

	repos := testtarget.ReposFlags{}
	if repos.Limit != 0 {
		t.Errorf("expected zero Limit, got %d", repos.Limit)
	}
	if repos.Format != "" {
		t.Errorf("expected empty Format, got %q", repos.Format)
	}

	list := testtarget.ListFlags{}
	if list.Limit != 0 {
		t.Errorf("expected zero Limit, got %d", list.Limit)
	}
	if list.Format != "" {
		t.Errorf("expected empty Format, got %q", list.Format)
	}

	fork := testtarget.ForkFlags{}
	if fork.To != "" {
		t.Errorf("expected empty To, got %q", fork.To)
	}
	if fork.RequestTimeout != 0 {
		t.Errorf("expected zero RequestTimeout, got %v", fork.RequestTimeout)
	}
}
