package app

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	config struct {
		Roots []string
	}
	configPath string
)

func DefaultRoot() string {
	return expandPath(config.Roots[0])
}

func Roots() []string {
	roots := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		roots = append(roots, expandPath(r))
	}
	return roots
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
