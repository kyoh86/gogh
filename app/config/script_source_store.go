package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/script"
)

type ScriptSourceStore struct{}

func NewScriptSourceStore() *ScriptSourceStore { return &ScriptSourceStore{} }

func (cs *ScriptSourceStore) Save(ctx context.Context, scriptID string, content io.Reader) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	if err := os.MkdirAll(source, 0755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}
	filePath := filepath.Join(source, scriptID)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create content file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("failed to write script: %w", err)
	}
	return nil
}

func (cs *ScriptSourceStore) Open(ctx context.Context, scriptID string) (io.ReadCloser, error) {
	source, err := cs.Source()
	if err != nil {
		return nil, fmt.Errorf("get content source: %w", err)
	}
	return os.Open(filepath.Join(source, scriptID))
}

func (cs *ScriptSourceStore) Remove(ctx context.Context, scriptID string) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	return os.Remove(filepath.Join(source, scriptID))
}

func (*ScriptSourceStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_HOOK_CONTENT_PATH", os.UserConfigDir, "script.v4")
	if err != nil {
		return "", fmt.Errorf("search script content path: %w", err)
	}
	return path, nil
}

var _ script.ScriptSourceStore = (*ScriptSourceStore)(nil)
