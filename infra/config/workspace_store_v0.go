package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	"gopkg.in/yaml.v2"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStoreV0 struct {
	filename string
}

type yamlWorkspaceStoreV0 struct {
	Roots []workspace.Root `yaml:"roots,omitempty"`
}

// Load implements workspace.WorkspaceRepository.
func (w *WorkspaceStoreV0) Load(ctx context.Context) (workspace.WorkspaceService, error) {
	var v yamlWorkspaceStoreV0
	file, err := os.Open(w.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := filesystem.WorkspaceService{}
	for i, root := range v.Roots {
		if err := svc.AddRoot(root, i == 0); err != nil {
			return nil, err
		}
	}
	svc.MarkSaved()
	return &svc, nil
}

func WorkspacePathV0() (string, error) {
	path, err := appContextPath("GOGH_CONFIG_PATH", os.UserConfigDir, "config.yaml")
	if err != nil {
		return "", fmt.Errorf("search config path: %w", err)
	}
	return path, nil
}

// NewWorkspaceStore creates a new WorkspaceStore instance.
func NewWorkspaceStoreV0(filename string) *WorkspaceStoreV0 {
	return &WorkspaceStoreV0{
		filename: filename,
	}
}

var _ store.Loader[workspace.WorkspaceService] = (*WorkspaceStoreV0)(nil)
