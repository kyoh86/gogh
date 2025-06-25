package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/store"
)

// ScriptDir returns the path to the script directory.
func ScriptDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_SCRIPT_PATH", os.UserConfigDir, "script.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search script path: %w", err)
	}
	return path, nil
}

type tomlScript struct {
	ID   uuid.UUID `toml:"id"`
	Name string    `toml:"name"`

	CreatedAt time.Time `toml:"created-at"`
	UpdatedAt time.Time `toml:"updated-at"`
}

type tomlScriptStore struct {
	Scripts []tomlScript `toml:"scripts"`
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

	data, err := loadTOMLFile[tomlScriptStore](src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("load script store: %w", err)
	}

	svc := initial()
	if err := svc.Load(func(yield func(script.Script, error) bool) {
		for _, s := range data.Scripts {
			if !yield(script.ConcreteScript(
				s.ID,
				s.Name,
				s.CreatedAt,
				s.UpdatedAt,
			), nil) {
				return
			}
		}
	}); err != nil {
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
		data.Scripts = append(data.Scripts, tomlScript{
			ID:        h.UUID(),
			Name:      h.Name(),
			CreatedAt: h.CreatedAt(),
			UpdatedAt: h.UpdatedAt(),
		})
	}
	if err := saveTOMLFile(src, data); err != nil {
		return fmt.Errorf("save script store: %w", err)
	}
	svc.MarkSaved()
	return nil
}

var _ store.Store[script.ScriptService] = (*ScriptStore)(nil)
