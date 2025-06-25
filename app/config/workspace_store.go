package config

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// WorkspaceStore is a repository for managing workspace configuration.
type WorkspaceStore struct{}

type tomlWorkspaceStore struct {
	Roots       []workspace.Root `toml:"roots,omitempty"`
	PrimaryRoot string           `toml:"primary_root,omitempty"`
}

// Load implements store.Store.
func (w *WorkspaceStore) Load(ctx context.Context, initial func() workspace.WorkspaceService) (workspace.WorkspaceService, error) {
	source, err := w.Source()
	if err != nil {
		return nil, err
	}

	v, err := loadTOMLFile[tomlWorkspaceStore](source)
	if err != nil {
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
	v := tomlWorkspaceStore{
		Roots:       ws.GetRoots(),
		PrimaryRoot: ws.GetPrimaryRoot(),
	}

	if err := saveTOMLFile(source, v); err != nil {
		return err
	}
	ws.MarkSaved()
	return nil
}

func (*WorkspaceStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_WORKSPACE_PATH", os.UserConfigDir, "workspace.v4.toml")
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
