package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/pelletier/go-toml/v2"
)

type OverlayStore struct{}

type tomlOverlayFile struct {
	SourcePath string `toml:"source_path"`
	TargetPath string `toml:"target_path"`
}

type tomlOverlayPattern struct {
	Pattern string            `toml:"pattern"`
	Files   []tomlOverlayFile `toml:"files"`
}

type tomlOverlayStore struct {
	Patterns []tomlOverlayPattern `toml:"patterns"`
}

// Load implements store.Store interface
func (s *OverlayStore) Load(ctx context.Context, initial func() workspace.OverlayService) (workspace.OverlayService, error) {
	var v tomlOverlayStore
	source, err := s.Source()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(source)
	if os.IsNotExist(err) {
		return initial(), nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := toml.NewDecoder(file).Decode(&v); err != nil {
		return nil, err
	}

	svc := initial()
	for _, pattern := range v.Patterns {
		files := make([]workspace.OverlayFile, 0, len(pattern.Files))
		for _, file := range pattern.Files {
			files = append(files, workspace.OverlayFile{
				SourcePath: file.SourcePath,
				TargetPath: file.TargetPath,
			})
		}
		if err := svc.AddPattern(pattern.Pattern, files); err != nil {
			return nil, fmt.Errorf("add pattern %s: %w", pattern.Pattern, err)
		}
	}
	svc.MarkSaved()
	return svc, nil
}

// Save implements store.Store interface
func (s *OverlayStore) Save(ctx context.Context, svc workspace.OverlayService, force bool) error {
	if !svc.HasChanges() && !force {
		return nil
	}
	source, err := s.Source()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	patterns := svc.GetPatterns()
	v := tomlOverlayStore{
		Patterns: make([]tomlOverlayPattern, 0, len(patterns)),
	}

	for _, pattern := range patterns {
		files := make([]tomlOverlayFile, 0, len(pattern.Files))
		for _, file := range pattern.Files {
			files = append(files, tomlOverlayFile{
				SourcePath: file.SourcePath,
				TargetPath: file.TargetPath,
			})
		}
		v.Patterns = append(v.Patterns, tomlOverlayPattern{
			Pattern: pattern.Pattern,
			Files:   files,
		})
	}

	if err := toml.NewEncoder(file).Encode(v); err != nil {
		return err
	}
	svc.MarkSaved()
	return nil
}

func (*OverlayStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_PATH", os.UserConfigDir, "overlay.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search overlay path: %w", err)
	}
	return path, nil
}

func NewOverlayStore() *OverlayStore {
	return &OverlayStore{}
}

var _ store.Saver[workspace.OverlayService] = (*OverlayStore)(nil)
