package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/pelletier/go-toml/v2"
)

// HookDir returns the path to the hook directory.
func HookDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_HOOK_PATH", os.UserConfigDir, "hook.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search hook path: %w", err)
	}
	return path, nil
}

type tomlHook struct {
	ID   uuid.UUID `toml:"id"`
	Name string    `toml:"name"`

	RepoPattern  string `toml:"repo-pattern"`
	TriggerEvent string `toml:"trigger-event"`

	OperationType string `toml:"operation-type"`
	OperationID   string `toml:"operation-id"`
}

type tomlHookStore struct {
	Hooks []tomlHook `toml:"hooks"`
}

type HookStore struct{}

func NewHookStore() *HookStore { return &HookStore{} }

func (s *HookStore) Source() (string, error) {
	return HookDir()
}

func (s *HookStore) Load(ctx context.Context, initial func() hook.HookService) (hook.HookService, error) {
	src, err := s.Source()
	if err != nil {
		return nil, fmt.Errorf("get hook store source: %w", err)
	}
	f, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("open hook store: %w", err)
	}
	defer f.Close()
	var data tomlHookStore
	if err := toml.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode hook store: %w", err)
	}
	svc := initial()
	if err := svc.Load(func(yield func(hook.Hook, error) bool) {
		for _, h := range data.Hooks {
			if !yield(hook.ConcreteHook(
				h.ID,
				h.Name,
				h.RepoPattern,
				h.TriggerEvent,
				h.OperationType,
				h.OperationID,
			), nil) {
				return
			}
		}
	}); err != nil {
		return nil, fmt.Errorf("set hooks: %w", err)
	}
	svc.MarkSaved()
	return svc, nil
}

func (s *HookStore) Save(ctx context.Context, svc hook.HookService, force bool) error {
	if !svc.HasChanges() && !force {
		return nil
	}
	src, err := s.Source()
	if err != nil {
		return fmt.Errorf("get hook store source: %w", err)
	}
	data := tomlHookStore{}
	for h, err := range svc.List() {
		if err != nil {
			return fmt.Errorf("list hooks: %w", err)
		}
		data.Hooks = append(data.Hooks, tomlHook{
			ID:            h.UUID(),
			Name:          h.Name(),
			RepoPattern:   h.RepoPattern(),
			TriggerEvent:  string(h.TriggerEvent()),
			OperationType: string(h.OperationType()),
			OperationID:   h.OperationID(),
		})
	}
	if err := os.MkdirAll(filepath.Dir(src), 0755); err != nil {
		return fmt.Errorf("create hook store directory: %w", err)
	}
	f, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open hook store file: %w", err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return fmt.Errorf("encode hook store file: %w", err)
	}
	svc.MarkSaved()
	return nil
}

var _ store.Store[hook.HookService] = (*HookStore)(nil)
