package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/pelletier/go-toml/v2"
)

// ScriptDir returns the path to the script directory.
func ScriptDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_SCRIPT_PATH", os.UserConfigDir, "script.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search script path: %w", err)
	}
	return path, nil
}

type tomlScriptStore struct {
	Scripts []script.Script `toml:"scripts"`
	// TODO: script.Script is not marshalable by toml
}

type ScriptStore struct{}

func NewScriptStore() *ScriptStore { return &ScriptStore{} }

func (s *ScriptStore) Source() (string, error) {
	return ScriptDir()
}

func (s *ScriptStore) Load(ctx context.Context, initial func() script.ScriptService) (script.ScriptService, error) {
	src, err := s.Source()
	if err != nil {
		return nil, fmt.Errorf("get script store source: %w", err)
	}
	f, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("open script store: %w", err)
	}
	defer f.Close()
	var data tomlScriptStore
	if err := toml.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode script store: %w", err)
	}
	svc := initial()
	if err := svc.Load(typ.WithNilError(slices.Values(data.Scripts))); err != nil {
		return nil, fmt.Errorf("set scripts: %w", err)
	}
	svc.MarkSaved()
	return svc, nil
}

func (s *ScriptStore) Save(ctx context.Context, svc script.ScriptService, force bool) error {
	if !svc.HasChanges() && !force {
		return nil
	}
	src, err := s.Source()
	if err != nil {
		return fmt.Errorf("get script store source: %w", err)
	}
	data := tomlScriptStore{}
	for h, err := range svc.List() {
		if err != nil {
			return fmt.Errorf("list scripts: %w", err)
		}
		data.Scripts = append(data.Scripts, h)
	}
	if err := os.MkdirAll(filepath.Dir(src), 0755); err != nil {
		return fmt.Errorf("create script store directory: %w", err)
	}
	f, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open script store file: %w", err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return fmt.Errorf("encode script store file: %w", err)
	}
	svc.MarkSaved()
	return nil
}

var _ store.Store[script.ScriptService] = (*ScriptStore)(nil)
