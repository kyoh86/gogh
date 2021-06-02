package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type FlagStruct struct {
	Create struct {
		Template            string
		LicenseTemplate     string
		GitignoreTemplate   string
		Private             bool
		DisableDownloads    bool
		DisableWiki         bool
		AutoInit            bool
		DisableProjects     bool
		DisableIssues       bool
		PreventSquashMerge  bool
		PreventMergeCommit  bool
		PreventRebaseMerge  bool
		DeleteBranchOnMerge bool
	}
}

var (
	config struct {
		Roots []string
		Flag  FlagStruct
	}
	configPath string
)

func DefaultRoot() string {
	return expandPath(config.Roots[0])
}

func Flag() FlagStruct {
	return config.Flag
}

func Roots() []string {
	roots := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		roots = append(roots, expandPath(r))
	}
	return roots
}

func SetDefaultRoot(r string) {
	roots := make([]string, 0, len(config.Roots))
	roots = append(roots, r)
	for _, root := range config.Roots {
		if root == r {
			continue
		}
		roots = append(roots, root)
	}
	config.Roots = roots
}

func AddRoots(roots []string) {
	config.Roots = append(config.Roots, roots...)
}

func RemoveRoot(r string) {
	roots := make([]string, 0, len(config.Roots))
	for _, root := range config.Roots {
		if root == r {
			continue
		}
		roots = append(roots, root)
	}
	config.Roots = roots
}

func expandPath(p string) string {
	p = os.ExpandEnv(p)
	runes := []rune(p)
	if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
		return filepath.Join(homeDir, string(runes[2:]))
	}
	return p
}

func setupConfig() error {
	configPath = filepath.Join(configDir, Name, "config.yaml")
	if err := loadYAML(configPath, &config); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if len(config.Roots) == 0 {
		config.Roots = []string{
			filepath.Join(homeDir, "Projects"),
		}
	}
	return nil
}

func SaveConfig() error {
	if err := saveYAML(configPath, config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}
