package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestLocationFormatter(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		wantErr   bool
		wantValue interface{}
	}{
		{
			name:      "empty string returns path format",
			input:     "",
			wantErr:   false,
			wantValue: repository.LocationFormatPath,
		},
		{
			name:      "rel-path returns path format",
			input:     "rel-path",
			wantErr:   false,
			wantValue: repository.LocationFormatPath,
		},
		{
			name:      "rel returns path format",
			input:     "rel",
			wantErr:   false,
			wantValue: repository.LocationFormatPath,
		},
		{
			name:      "path returns path format",
			input:     "path",
			wantErr:   false,
			wantValue: repository.LocationFormatPath,
		},
		{
			name:      "rel-file-path returns path format",
			input:     "rel-file-path",
			wantErr:   false,
			wantValue: repository.LocationFormatPath,
		},
		{
			name:      "full-file-path returns full path format",
			input:     "full-file-path",
			wantErr:   false,
			wantValue: repository.LocationFormatFullPath,
		},
		{
			name:      "full returns full path format",
			input:     "full",
			wantErr:   false,
			wantValue: repository.LocationFormatFullPath,
		},
		{
			name:      "json returns JSON format",
			input:     "json",
			wantErr:   false,
			wantValue: repository.LocationFormatJSON,
		},
		{
			name:      "fields returns fields format with tab separator",
			input:     "fields",
			wantErr:   false,
			wantValue: repository.LocationFormatFields("\t"),
		},
		{
			name:      "fields with custom separator",
			input:     "fields:|",
			wantErr:   false,
			wantValue: repository.LocationFormatFields("|"),
		},
		{
			name:      "fields with comma separator",
			input:     "fields:,",
			wantErr:   false,
			wantValue: repository.LocationFormatFields(","),
		},
		{
			name:      "fields with space separator",
			input:     "fields: ",
			wantErr:   false,
			wantValue: repository.LocationFormatFields(" "),
		},
		{
			name:      "invalid format",
			input:     "invalid",
			wantErr:   true,
			wantValue: nil,
		},
		{
			name:      "unknown format",
			input:     "unknown-format",
			wantErr:   true,
			wantValue: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := testtarget.LocationFormatter(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("LocationFormatter(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				// Compare the formatters by their string representation
				gotStr := formatToString(got)
				wantStr := formatToString(tc.wantValue.(repository.LocationFormat))
				if gotStr != wantStr {
					t.Errorf("LocationFormatter(%q) = %v, want %v", tc.input, gotStr, wantStr)
				}
			}
		})
	}
}

// Helper function to convert LocationFormat to string for comparison
func formatToString(f repository.LocationFormat) string {
	if f == nil {
		return "nil"
	}
	// Test with a sample location
	loc := repository.NewLocation("/path/to/repo", "github.com", "owner", "name")
	formatted, err := f.Format(*loc)
	if err != nil {
		return "error: " + err.Error()
	}
	return formatted
}

func TestFlags_HasChanges(t *testing.T) {
	f := &testtarget.Flags{}
	if f.HasChanges() {
		t.Errorf("HasChanges should always return false")
	}

	// Test with RawHasChanges set to true
	f.RawHasChanges = true
	if !f.HasChanges() {
		t.Errorf("HasChanges should return true when RawHasChanges is true")
	}
}

func TestFlags_MarkSaved(t *testing.T) {
	f := &testtarget.Flags{}
	// Set RawHasChanges to true
	f.RawHasChanges = true
	f.MarkSaved()
	// Verify it's set back to false
	if f.RawHasChanges {
		t.Errorf("RawHasChanges should be false after MarkSaved")
	}
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
	if f.BundleRestore.CloneRetryTimeout != 5*time.Minute {
		t.Errorf("expected BundleRestore.CloneRetryTimeout to be 5m, got %v", f.BundleRestore.CloneRetryTimeout)
	}
	if f.Clone.CloneRetryTimeout != 5*time.Minute {
		t.Errorf("expected Clone.CloneRetryTimeout to be 5m, got %v", f.Clone.CloneRetryTimeout)
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
	if f.Create.CloneRetryTimeout != 5*time.Minute {
		t.Errorf("expected Create.CloneRetryTimeout to be 5m, got %v", f.Create.CloneRetryTimeout)
	}
	if f.Create.CloneRetryLimit != 3 {
		t.Errorf("expected Create.CloneRetryLimit to be 3, got %d", f.Create.CloneRetryLimit)
	}
	if f.List.Limit != 100 {
		t.Errorf("expected List.Limit to be 100, got %d", f.List.Limit)
	}
	if f.Fork.CloneRetryTimeout != 5*time.Minute {
		t.Errorf("expected Fork.CloneRetryTimeout to be 5m, got %v", f.Fork.CloneRetryTimeout)
	}
	if f.Fork.CloneRetryLimit != 3 {
		t.Errorf("expected Fork.CloneRetryLimit to be 3, got %d", f.Fork.CloneRetryLimit)
	}

	// Default boolean flags should be false
	if f.BundleRestore.DryRun {
		t.Errorf("expected BundleRestore.DryRun to be false")
	}
	if f.Clone.DryRun {
		t.Errorf("expected Clone.DryRun to be false")
	}
	if f.Create.DryRun {
		t.Errorf("expected Create.DryRun to be false")
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
	if bundleRestore.CloneRetryTimeout != 0 {
		t.Errorf("expected zero CloneRetryTimeout, got %v", bundleRestore.CloneRetryTimeout)
	}
	if bundleRestore.File != "" {
		t.Errorf("expected empty File, got %q", bundleRestore.File)
	}
	if bundleRestore.CloneRetryLimit != 0 {
		t.Errorf("expected zero CloneRetryLimit, got %d", bundleRestore.CloneRetryLimit)
	}
	if bundleRestore.DryRun {
		t.Errorf("expected false DryRun")
	}

	clone := testtarget.CloneFlags{}
	if clone.CloneRetryTimeout != 0 {
		t.Errorf("expected zero CloneRetryTimeout, got %v", clone.CloneRetryTimeout)
	}
	if clone.DryRun {
		t.Errorf("expected false DryRun")
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
	if fork.CloneRetryTimeout != 0 {
		t.Errorf("expected zero CloneRetryTimeout, got %v", fork.CloneRetryTimeout)
	}
}
