package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var config struct {
	Roots       []ExpandablePath `yaml:"roots"`
	DefaultFlag struct {
		Create createFlagsStruct `yaml:"create,omitempty"`
		Repos  reposFlagsStruct  `yaml:"repos,omitempty"`
	} `yaml:"flag,omitempty"`
}

func defaultRoot() string {
	return config.Roots[0].expanded
}

func Roots() []string {
	list := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		list = append(list, r.expanded)
	}
	return list
}

func setDefaultRoot(r string) error {
	roots := make([]ExpandablePath, 0, len(config.Roots))
	newDefault, err := ParsePath(r)
	if err != nil {
		return err
	}
	roots = append(roots, newDefault)
	for _, root := range config.Roots {
		if root.raw == r {
			continue
		}
		roots = append(roots, root)
	}
	config.Roots = roots
	return nil
}

func addRoots(roots []string) error {
	for _, r := range roots {
		newRoot, err := ParsePath(r)
		if err != nil {
			return err
		}
		config.Roots = append(config.Roots, newRoot)
	}
	return nil
}

func removeRoot(r string) {
	roots := make([]ExpandablePath, 0, len(config.Roots))
	for _, root := range config.Roots {
		if root.raw == r {
			continue
		}
		roots = append(roots, root)
	}
	config.Roots = roots
}

func loadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("search user config dir: %w", err)
	}
	configPath := filepath.Join(configDir, appName, "config.yaml")
	if err := loadYAML(configPath, &config); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if len(config.Roots) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("search user home dir: %w", err)
		}
		raw := filepath.Join(homeDir, "Projects")
		config.Roots = []ExpandablePath{{
			raw:      raw,
			expanded: raw,
		}}
	}
	return nil
}

func SaveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("search user config dir: %w", err)
	}
	configPath := filepath.Join(configDir, appName, "config.yaml")
	if err := saveYAML(configPath, config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}
