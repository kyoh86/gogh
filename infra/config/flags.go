package config

import (
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
)

type BundleRestoreFlags struct {
	File            Path `yaml:"file,omitempty" toml:"file,omitempty"`
	CloneRetryLimit int  `yaml:"cloneRetryLimit,omitempty" toml:"cloneRetryLimit,omitempty"`
	Dryrun          bool `yaml:"-" toml:"-"`
}

type ReposFlags struct {
	Limit       int                    `yaml:"limit,omitempty" toml:"limit,omitempty"`
	Public      bool                   `yaml:"public,omitempty" toml:"public,omitempty"`
	Private     bool                   `yaml:"private,omitempty" toml:"private,omitempty"`
	Fork        bool                   `yaml:"fork,omitempty" toml:"fork,omitempty"`
	NotFork     bool                   `yaml:"notFork,omitempty" toml:"notFork,omitempty"`
	Archived    bool                   `yaml:"archived,omitempty" toml:"archived,omitempty"`
	NotArchived bool                   `yaml:"notArchived,omitempty" toml:"notArchived,omitempty"`
	Format      flags.RemoteRepoFormat `yaml:"format,omitempty" toml:"format,omitempty"`
	Color       string                 `yaml:"color,omitempty" toml:"color,omitempty"`
	Relation    []string               `yaml:"relation,omitempty" toml:"relation,omitempty"`
	Sort        string                 `yaml:"sort,omitempty" toml:"sort,omitempty"`
	Order       string                 `yaml:"order,omitempty" toml:"order,omitempty"`
}

type CreateFlags struct {
	Template            string `yaml:"template,omitempty" toml:"template,omitempty"`
	Description         string `yaml:"-" toml:"-"`
	Homepage            string `yaml:"-" toml:"-"`
	LicenseTemplate     string `yaml:"licenseTemplate,omitempty" toml:"licenseTemplate,omitempty"`
	GitignoreTemplate   string `yaml:"gitignoreTemplate,omitempty" toml:"gitignoreTemplate,omitempty"`
	CloneRetryLimit     int    `yaml:"cloneRetryLimit,omitempty" toml:"cloneRetryLimit,omitempty"`
	DisableWiki         bool   `yaml:"disableWiki,omitempty" toml:"disableWiki,omitempty"`
	DisableDownloads    bool   `yaml:"disableDownloads,omitempty" toml:"disableDownloads,omitempty"`
	IsTemplate          bool   `yaml:"-" toml:"-"`
	AutoInit            bool   `yaml:"autoInit,omitempty" toml:"autoInit,omitempty"`
	DisableProjects     bool   `yaml:"disableProjects,omitempty" toml:"disableProjects,omitempty"`
	DisableIssues       bool   `yaml:"disableIssues,omitempty" toml:"disableIssues,omitempty"`
	PreventSquashMerge  bool   `yaml:"preventSquashMerge,omitempty" toml:"preventSquashMerge,omitempty"`
	PreventMergeCommit  bool   `yaml:"preventMergeCommit,omitempty" toml:"preventMergeCommit,omitempty"`
	PreventRebaseMerge  bool   `yaml:"preventRebaseMerge,omitempty" toml:"preventRebaseMerge,omitempty"`
	DeleteBranchOnMerge bool   `yaml:"deleteBranchOnMerge,omitempty" toml:"deleteBranchOnMerge,omitempty"`
	Private             bool   `yaml:"private,omitempty" toml:"private,omitempty"`
	IncludeAllBranches  bool   `yaml:"includeAllBranches,omitempty" toml:"includeAllBranches,omitempty"`
	Dryrun              bool   `yaml:"-" toml:"-"`
}

type CwdFlags struct {
	Format flags.LocalRepoFormat `yaml:"format,omitempty" toml:"format,omitempty"`
}

type ListFlags struct {
	Limit   int                   `yaml:"limit,omitempty" toml:"limit,omitempty"`
	Query   string                `yaml:"-" toml:"-"`
	Format  flags.LocalRepoFormat `yaml:"format,omitempty" toml:"format,omitempty"`
	Primary bool                  `yaml:"primary,omitempty" toml:"primary,omitempty"`
}

type ForkFlags struct {
	To                string `yaml:"to,omitempty" toml:"to,omitempty"`
	DefaultBranchOnly bool   `yaml:"defaultBranchOnly,omitempty" toml:"defaultBranchOnly,omitempty"`
	CloneRetryLimit   int    `yaml:"cloneRetryLimit,omitempty" toml:"cloneRetryLimit,omitempty"`
}

type BundleDumpFlags struct {
	File Path `yaml:"file,omitempty" toml:"file,omitempty"`
}

type Flags struct {
	BundleRestore BundleRestoreFlags `yaml:"bundleRestore,omitempty" toml:"bundleRestore,omitempty"`
	BundleDump    BundleDumpFlags    `yaml:"bundleDump,omitempty" toml:"bundleDump,omitempty"`
	List          ListFlags          `yaml:"list,omitempty" toml:"list,omitempty"`
	Cwd           CwdFlags           `yaml:"cwd,omitempty" toml:"cwd,omitempty"`
	Create        CreateFlags        `yaml:"create,omitempty" toml:"create,omitempty"`
	Repos         ReposFlags         `yaml:"repos,omitempty" toml:"repos,omitempty"`
	Fork          ForkFlags          `yaml:"fork,omitempty" toml:"fork,omitempty"`
}
