package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/pelletier/go-toml/v2"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStore struct{}

type tomlWorkspaceStore struct {
	Roots       []workspace.Root `toml:"roots,omitempty"`
	PrimaryRoot string           `toml:"primary_root,omitempty"`
}

// Load implements store.Store.
func (w *WorkspaceStore) Load(ctx context.Context, initial func() workspace.WorkspaceService) (workspace.WorkspaceService, error) {
	var v tomlWorkspaceStore
	source, err := w.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}
	svc := initial()
	for _, root := range v.Roots {
		if err := svc.AddRoot(root, root == v.PrimaryRoot); err != nil {
			return nil, err
		}
	}
	svc.MarkSaved()
	return svc, nil
}

// Save implements workspace.WorkspaceRepository.
func (w *WorkspaceStore) Save(ctx context.Context, ws workspace.WorkspaceService, force bool) error {
	if !ws.HasChanges() && !force {
		return nil
	}
	source, err := w.Source()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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
	ws.MarkSaved()
	return nil
}

func (*WorkspaceStore) Source() (string, error) {
	path, err := appContextPath("GOGH_WORKSPACE_PATH", os.UserConfigDir, "workspace.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search workspace path: %w", err)
	}
	return path, nil
}

// NewWorkspaceStore creates a new WorkspaceStore instance.
func NewWorkspaceStore() *WorkspaceStore {
	return &WorkspaceStore{}
}

var _ store.Store[workspace.WorkspaceService] = (*WorkspaceStore)(nil)
