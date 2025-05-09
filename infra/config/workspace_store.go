package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	"github.com/pelletier/go-toml/v2"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStore struct {
	filename string
}

type tomlWorkspaceStore struct {
	Roots       []workspace.Root `toml:"roots,omitempty"`
	PrimaryRoot string           `toml:"primary_root,omitempty"`
}

// Load implements workspace.WorkspaceRepository.
func (w *WorkspaceStore) Load(ctx context.Context) (workspace.WorkspaceService, error) {
	var v tomlWorkspaceStore
	file, err := os.Open(w.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := filesystem.NewWorkspaceService()
	for _, root := range v.Roots {
		if err := svc.AddRoot(root, root == v.PrimaryRoot); err != nil {
			return nil, err
		}
	}
	return svc, nil
}

// Save implements workspace.WorkspaceRepository.
func (w *WorkspaceStore) Save(ctx context.Context, ws workspace.WorkspaceService) error {
	file, err := os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	v := tomlWorkspaceStore{
		Roots:       ws.GetRoots(),
		PrimaryRoot: ws.GetPrimaryRoot(),
	}

	if err := toml.NewEncoder(file).Encode(v); err != nil {
		return err
	}
	return nil
}

func WorkspacePath() (string, error) {
	path, err := appContextPath("GOGH_WORKSPACE_PATH", os.UserCacheDir, "workspace.v1.yaml")
	if err != nil {
		return "", fmt.Errorf("search workspace path: %w", err)
	}
	return path, nil
}

// NewWorkspaceStore creates a new WorkspaceStore instance.
func NewWorkspaceStore(filename string) *WorkspaceStore {
	return &WorkspaceStore{
		filename: filename,
	}
}

var _ workspace.WorkspaceStore = (*WorkspaceStore)(nil)
