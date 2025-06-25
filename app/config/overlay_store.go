package config

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/store"
)

// OverlayDir returns the path to the overlay directory.
func OverlayDir() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_PATH", os.UserConfigDir, "overlay.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search overlay path: %w", err)
	}
	return path, nil
}

// tomlOverlay is used for (un)marshaling overlays to/from TOML.
type tomlOverlay struct {
	ID           uuid.UUID `toml:"id"`
	Name         string    `toml:"name"`
	RelativePath string    `toml:"relative-path"`
}

// tomlOverlayStore is used for (un)marshaling overlays to/from TOML.
type tomlOverlayStore struct {
	Overlays []tomlOverlay `toml:"overlays"`
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

	data, err := loadTOMLFile[tomlOverlayStore](src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("load overlay store: %w", err)
	}

	svc := initial()
	if err := svc.Load(func(yield func(overlay overlay.Overlay, err error) bool) {
		if err != nil {
			return
		}
		for _, o := range data.Overlays {
			if !yield(overlay.ConcreteOverlay(o.ID, o.Name, o.RelativePath), nil) {
				return
			}
		}
	}); err != nil {
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
		data.Overlays = append(data.Overlays, tomlOverlay{
			ID:           ov.UUID(),
			Name:         ov.Name(),
			RelativePath: ov.RelativePath(),
		})
	}

	if err := saveTOMLFile(src, data); err != nil {
		return fmt.Errorf("save overlay store: %w", err)
	}
	svc.MarkSaved()
	return nil
}

var _ store.Store[overlay.OverlayService] = (*OverlayStore)(nil)
