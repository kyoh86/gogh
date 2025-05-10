package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type ConfigStore struct {
	Roots []Path `yaml:"roots"`
}

var (
	globalConfig ConfigStore
	configOnce   sync.Once
)

func ConfigPath() (string, error) {
	path, err := appContextPath("GOGH_CONFIG_PATH", os.UserConfigDir, "config.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

func LoadConfig() (_ *ConfigStore, retErr error) {
	configOnce.Do(func() {
		path, err := ConfigPath()
		if err != nil {
			retErr = err
			return
		}

		if err := loadYAML(path, &globalConfig); err != nil {
			retErr = err
			return
		}
		if len(globalConfig.Roots) == 0 {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				retErr = fmt.Errorf("search user home dir: %w", err)
				return
			}
			raw := filepath.Join(homeDir, "Projects")
			globalConfig.Roots = []Path{{
				raw:      raw,
				expanded: raw,
			}}
		}
	})
	return &globalConfig, retErr
}

func SaveConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return saveYAML(path, globalConfig)
}

func (c *ConfigStore) PrimaryRoot() string {
	return c.Roots[0].expanded
}

func (c *ConfigStore) GetRoots() []string {
	list := make([]string, 0, len(c.Roots))
	for _, r := range c.Roots {
		list = append(list, r.expanded)
	}
	return list
}

func (c *ConfigStore) SetPrimaryRoot(r string) error {
	rootList := make([]Path, 0, len(c.Roots))
	newDefault, err := parsePath(r)
	if err != nil {
		return err
	}
	rootList = append(rootList, newDefault)
	for _, root := range c.Roots {
		if root.raw == r {
			continue
		}
		rootList = append(rootList, root)
	}
	globalConfig.Roots = rootList
	return nil
}

func (c *ConfigStore) AddRoots(rootList []string) error {
	for _, r := range rootList {
		newRoot, err := parsePath(r)
		if err != nil {
			return err
		}
		c.Roots = append(c.Roots, newRoot)
	}
	return nil
}

func (c *ConfigStore) RemoveRoot(r string) {
	rootList := make([]Path, 0, len(c.Roots))
	for _, root := range c.Roots {
		if root.raw == r || root.expanded == r {
			continue
		}
		rootList = append(rootList, root)
	}
	c.Roots = rootList
}
