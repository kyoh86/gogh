package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"gopkg.in/yaml.v2"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStoreV0 struct{}

type yamlWorkspaceStoreV0 struct {
	Roots []workspace.Root `yaml:"roots,omitempty"`
}

func ExpandPath(p string) (string, error) {
	p = os.ExpandEnv(p)
	runes := []rune(p)
	switch len(runes) {
	case 0:
		return p, nil
	case 1:
		if runes[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("search user home dir: %w", err)
			}
			return homeDir, nil
		}
	default:
		if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("search user home dir: %w", err)
			}
			return filepath.Join(homeDir, string(runes[2:])), nil
		}
	}
	return p, nil
}

// Load implements store.Loader
func (w *WorkspaceStoreV0) Load(ctx context.Context, initial func() workspace.WorkspaceService) (workspace.WorkspaceService, error) {
	var v yamlWorkspaceStoreV0
	source, err := w.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := initial()
	for i, root := range v.Roots {
		path, err := ExpandPath(root)
		if err != nil {
			return nil, fmt.Errorf("expand path: %w", err)
		}
		if err := svc.AddRoot(path, i == 0); err != nil {
			return nil, err
		}
	}
	svc.MarkSaved()
	return svc, nil
}

func (*WorkspaceStoreV0) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_CONFIG_PATH", os.UserConfigDir, "config.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

// NewWorkspaceStore creates a new WorkspaceStore instance.
func NewWorkspaceStoreV0() *WorkspaceStoreV0 {
	return &WorkspaceStoreV0{}
}

var _ store.Loader[workspace.WorkspaceService] = (*WorkspaceStoreV0)(nil)
