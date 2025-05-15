package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v3/core/repository"
)

func LocationFormatter(v string) (repository.LocationFormat, error) {
	switch v {
	case "", "rel-path", "rel", "path", "rel-file-path":
		return repository.LocationFormatPath, nil
	case "full-file-path", "full":
		return repository.LocationFormatFullPath, nil
	case "json":
		return repository.LocationFormatJSON, nil
	case "fields":
		return repository.LocationFormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return repository.LocationFormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

// BundleDumpFlags is a struct that contains flags for dumping a bundle.
type BundleDumpFlags struct {
	File string `yaml:"file,omitempty" toml:"file,omitempty"`
}

// BundleRestoreFlags is a struct that contains flags for restoring a bundle.
type BundleRestoreFlags struct {
	File            string `yaml:"file,omitempty" toml:"file,omitempty"`
	CloneRetryLimit int    `yaml:"cloneRetryLimit,omitempty" toml:"clone-retry-limit,omitempty"`
	Dryrun          bool   `yaml:"-" toml:"-"`
}

// CreateFlags is a struct that contains flags for creating a repository.
type CreateFlags struct {
	Template            string `yaml:"template,omitempty" toml:"template,omitempty"`
	Description         string `yaml:"-" toml:"-"`
	Homepage            string `yaml:"-" toml:"-"`
	LicenseTemplate     string `yaml:"licenseTemplate,omitempty" toml:"license-template,omitempty"`
	GitignoreTemplate   string `yaml:"gitignoreTemplate,omitempty" toml:"gitignore-template,omitempty"`
	CloneRetryLimit     int    `yaml:"cloneRetryLimit,omitempty" toml:"clone-retry-limit,omitempty"`
	DisableWiki         bool   `yaml:"disableWiki,omitempty" toml:"disable-wiki,omitempty"`
	DisableDownloads    bool   `yaml:"disableDownloads,omitempty" toml:"disable-downloads,omitempty"`
	IsTemplate          bool   `yaml:"-" toml:"-"`
	AutoInit            bool   `yaml:"autoInit,omitempty" toml:"auto-init,omitempty"`
	DisableProjects     bool   `yaml:"disableProjects,omitempty" toml:"disable-projects,omitempty"`
	DisableIssues       bool   `yaml:"disableIssues,omitempty" toml:"disable-issues,omitempty"`
	PreventSquashMerge  bool   `yaml:"preventSquashMerge,omitempty" toml:"prevent-squash-merge,omitempty"`
	PreventMergeCommit  bool   `yaml:"preventMergeCommit,omitempty" toml:"prevent-merge-commit,omitempty"`
	PreventRebaseMerge  bool   `yaml:"preventRebaseMerge,omitempty" toml:"prevent-rebase-merge,omitempty"`
	DeleteBranchOnMerge bool   `yaml:"deleteBranchOnMerge,omitempty" toml:"delete-branch-on-merge,omitempty"`
	Private             bool   `yaml:"private,omitempty" toml:"private,omitempty"`
	IncludeAllBranches  bool   `yaml:"includeAllBranches,omitempty" toml:"include-all-branches,omitempty"`
	Dryrun              bool   `yaml:"-" toml:"-"`
}

// CwdFlags is a struct that contains flags for the cwd command.
type CwdFlags struct {
	Format string `yaml:"format,omitempty" toml:"format,omitempty"`
}

// ReposFlags is a struct that contains flags for the repos command.
type ReposFlags struct {
	Limit    int      `yaml:"limit,omitempty" toml:"limit,omitempty"`
	Privacy  string   `yaml:"privacy,omitempty" toml:"privacy,omitempty"`
	Fork     string   `yaml:"fork,omitempty" toml:"fork,omitempty"`
	Archive  string   `yaml:"archived,omitempty" toml:"archived,omitempty"`
	Format   string   `yaml:"format,omitempty" toml:"format,omitempty"`
	Color    string   `yaml:"color,omitempty" toml:"color,omitempty"`
	Relation []string `yaml:"relation,omitempty" toml:"relation,omitempty"`
	Sort     string   `yaml:"sort,omitempty" toml:"sort,omitempty"`
	Order    string   `yaml:"order,omitempty" toml:"order,omitempty"`
}

// ListFlags is a struct that contains flags for listing repositories.
type ListFlags struct {
	Limit   int    `yaml:"limit,omitempty" toml:"limit,omitempty"`
	Query   string `yaml:"-" toml:"-"`
	Format  string `yaml:"format,omitempty" toml:"format,omitempty"`
	Primary bool   `yaml:"primary,omitempty" toml:"primary,omitempty"`
}

// ForkFlags is a struct that contains flags for forking a repository.
type ForkFlags struct {
	To                string `yaml:"-" toml:"-"`
	DefaultBranchOnly bool   `yaml:"defaultBranchOnly,omitempty" toml:"default-branch-only,omitempty"`
	CloneRetryLimit   int    `yaml:"cloneRetryLimit,omitempty" toml:"clone-retry-limit,omitempty"`
}

// Flags is a struct that contains all the flags for the application.
type Flags struct {
	BundleDump    BundleDumpFlags    `yaml:"bundleDump,omitempty" toml:"bundle-dump,omitempty"`
	BundleRestore BundleRestoreFlags `yaml:"bundleRestore,omitempty" toml:"bundle-restore,omitempty"`
	List          ListFlags          `yaml:"list,omitempty" toml:"list,omitempty"`
	Cwd           CwdFlags           `yaml:"cwd,omitempty" toml:"cwd,omitempty"`
	Create        CreateFlags        `yaml:"create,omitempty" toml:"create,omitempty"`
	Repos         ReposFlags         `yaml:"repos,omitempty" toml:"repos,omitempty"`
	Fork          ForkFlags          `yaml:"fork,omitempty" toml:"fork,omitempty"`
}

// HasChanges always returns false because Flags does not support saving.
func (f *Flags) HasChanges() bool {
	return false
}

// MarkSaved is a no-op function. It does not save any changes to the Flags struct.
func (f *Flags) MarkSaved() {
	// No-op
}

func DefaultFlags() *Flags {
	f := new(Flags)
	homeDir, err := os.UserHomeDir()
	if err == nil && homeDir != "" {
		f.BundleDump.File = filepath.Join(homeDir, "./.config/gogh/bundle.txt")
		f.BundleRestore.File = filepath.Join(homeDir, "./.config/gogh/bundle.txt")
	}
	f.BundleRestore.CloneRetryLimit = 3

	f.Repos.Limit = 30
	f.Repos.Color = "auto"
	f.Repos.Relation = []string{"owner", "organization-member"}

	f.Create.CloneRetryLimit = 3

	f.List.Limit = 100

	f.Fork.CloneRetryLimit = 3
	return f
}
