package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/kyoh86/gogh/v3/ui/cli/flags"
)

type BundleRestoreFlags struct {
	File   Path `yaml:"file,omitempty"`
	Dryrun bool `yaml:"-"`
}

type ReposFlags struct {
	Format      string   `yaml:"format,omitempty"`
	Color       string   `yaml:"color,omitempty"`
	Sort        string   `yaml:"sort,omitempty"`
	Order       string   `yaml:"order,omitempty"`
	Relation    []string `yaml:"relation,omitempty"`
	Limit       int      `yaml:"limit,omitempty"`
	Private     bool     `yaml:"private,omitempty"`
	Public      bool     `yaml:"public,omitempty"`
	Fork        bool     `yaml:"fork,omitempty"`
	NotFork     bool     `yaml:"notFork,omitempty"`
	Archived    bool     `yaml:"archived,omitempty"`
	NotArchived bool     `yaml:"notArchived,omitempty"`
}

type CreateFlags struct {
	Template            string `yaml:"template,omitempty"`
	Description         string `yaml:"-"`
	Homepage            string `yaml:"-"`
	LicenseTemplate     string `yaml:"licenseTemplate,omitempty"`
	GitignoreTemplate   string `yaml:"gitignoreTemplate,omitempty"`
	CloneRetryLimit     int    `yaml:"cloneRetryLimit,omitempty"`
	DisableWiki         bool   `yaml:"disableWiki,omitempty"`
	DisableDownloads    bool   `yaml:"disableDownloads,omitempty"`
	IsTemplate          bool   `yaml:"-"`
	AutoInit            bool   `yaml:"autoInit,omitempty"`
	DisableProjects     bool   `yaml:"disableProjects,omitempty"`
	DisableIssues       bool   `yaml:"disableIssues,omitempty"`
	PreventSquashMerge  bool   `yaml:"preventSquashMerge,omitempty"`
	PreventMergeCommit  bool   `yaml:"preventMergeCommit,omitempty"`
	PreventRebaseMerge  bool   `yaml:"preventRebaseMerge,omitempty"`
	DeleteBranchOnMerge bool   `yaml:"deleteBranchOnMerge,omitempty"`
	Private             bool   `yaml:"private,omitempty"`
	Dryrun              bool   `yaml:"-"`
}

type CwdFlags struct {
	Format flags.ProjectFormat `yaml:"format,omitempty"`
}
type ListFlags struct {
	Query   string              `yaml:"-"`
	Format  flags.ProjectFormat `yaml:"format,omitempty"`
	Primary bool                `yaml:"primary,omitempty"`
}

type ForkFlags struct {
	Own bool `yaml:"own,omitempty"`
}

type BundleDumpFlags struct {
	File Path `yaml:"file,omitempty"`
}

type Flags struct {
	BundleRestore BundleRestoreFlags `yaml:"bundleRestore,omitempty"`
	BundleDump    BundleDumpFlags    `yaml:"bundleDump,omitempty"`
	List          ListFlags          `yaml:"list,omitempty"`
	Cwd           CwdFlags           `yaml:"cwd,omitempty"`
	Create        CreateFlags        `yaml:"create,omitempty"`
	Repos         ReposFlags         `yaml:"repos,omitempty"`
	Fork          ForkFlags          `yaml:"fork,omitempty"`
}

var (
	globalFlags Flags
	flagsOnce   sync.Once
)

func FlagsPath() (string, error) {
	path, err := appFilePath("GOGH_FLAG_PATH", os.UserConfigDir, "flag.yaml")
	if err != nil {
		return "", fmt.Errorf("search flags path: %w", err)
	}
	return path, nil
}

func LoadFlags() (_ *Flags, retErr error) {
	flagsOnce.Do(func() {
		path, err := FlagsPath()
		if err != nil {
			retErr = err
			return
		}

		if err := loadYAML(path, &globalFlags); err != nil {
			retErr = err
			return
		}
	})
	return &globalFlags, retErr
}

func SaveFlags() error {
	path, err := FlagsPath()
	if err != nil {
		return err
	}
	return saveYAML(path, globalFlags)
}
