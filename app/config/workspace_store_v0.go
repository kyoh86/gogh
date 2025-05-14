package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"gopkg.in/yaml.v2"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStoreV0 struct{}

type yamlWorkspaceStoreV0 struct {
	Roots []workspace.Root `yaml:"roots,omitempty"`
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
		if err := svc.AddRoot(root, i == 0); err != nil {
			return nil, err
		}
	}
	svc.MarkSaved()
	return svc, nil
}

func (*WorkspaceStoreV0) Source() (string, error) {
	path, err := appContextPath("GOGH_CONFIG_PATH", os.UserConfigDir, "config.yaml")
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
