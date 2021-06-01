package app

import (
	"fmt"
	"path/filepath"

	"github.com/kyoh86/gogh/v2"
)

var (
	servers     gogh.Servers
	serversPath string
)

func Servers() *gogh.Servers {
	return &servers
}

func setupServers() error {
	serversPath = filepath.Join(cacheDir, Name, "servers.yaml")
	if err := loadYAML(serversPath, &servers); err != nil {
		return fmt.Errorf("load servers: %w", err)
	}
	return nil
}

func SaveServers() error {
	if err := saveYAML(serversPath, &servers); err != nil {
		return fmt.Errorf("save servers: %w", err)
	}
	return nil
}
