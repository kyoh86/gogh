package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/pelletier/go-toml/v2"
)

// OverlayDir returns the path to the overlay directory.
func OverlayDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_PATH", os.UserConfigDir, "overlay.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search overlay path: %w", err)
	}
	return path, nil
}

// tomlOverlayStore is used for (un)marshaling overlays to/from TOML.
type tomlOverlayStore struct {
	Overlays []overlay.Overlay `toml:"overlays"`
}

// OverlayStore implements overlay.OverlayStore and persists overlays as TOML.
type OverlayStore struct{}

func NewOverlayStore() *OverlayStore {
	return &OverlayStore{}
}

func (s *OverlayStore) Source() (string, error) {
	return OverlayDir()
}

func (s *OverlayStore) Load(ctx context.Context, initial func() overlay.OverlayService) (overlay.OverlayService, error) {
	src, err := s.Source()
	if err != nil {
		return nil, fmt.Errorf("get overlay store source: %w", err)
	}
	f, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("open overlay store: %w", err)
	}
	defer f.Close()
	var data tomlOverlayStore
	if err := toml.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode overlay store: %w", err)
	}
	svc := initial()
	if err := svc.Set(data.Overlays); err != nil {
		return nil, fmt.Errorf("set overlays: %w", err)
	}
	svc.MarkSaved()
	return svc, nil
}

func (s *OverlayStore) Save(ctx context.Context, svc overlay.OverlayService, force bool) error {
	if !svc.HasChanges() && !force {
		return nil
	}
	src, err := s.Source()
	if err != nil {
		return fmt.Errorf("get overlay store source: %w", err)
	}
	data := tomlOverlayStore{}
	for ov, err := range svc.List() {
		if err != nil {
			return fmt.Errorf("list overlays: %w", err)
		}
		if ov == nil {
			continue
		}
		data.Overlays = append(data.Overlays, *ov)
	}
	if err := os.MkdirAll(filepath.Dir(src), 0755); err != nil {
		return fmt.Errorf("create overlay store directory: %w", err)
	}
	f, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open overlay store file: %w", err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return fmt.Errorf("encode overlay store file: %w", err)
	}
	svc.MarkSaved()
	return nil
}

var _ overlay.OverlayStore = (*OverlayStore)(nil)
