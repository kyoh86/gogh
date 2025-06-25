package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

type OverlayContentStore struct{}

func NewOverlayContentStore() *OverlayContentStore {
	return &OverlayContentStore{}
}

func (cs *OverlayContentStore) Save(ctx context.Context, overlayID string, content io.Reader) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	if err := os.MkdirAll(source, 0o755); err != nil {
		return fmt.Errorf("failed to create content directory: %w", err)
	}
	filePath := filepath.Join(source, overlayID)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create content file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	return nil
}

func (cs *OverlayContentStore) Open(ctx context.Context, overlayID string) (io.ReadCloser, error) {
	source, err := cs.Source()
	if err != nil {
		return nil, fmt.Errorf("get content source: %w", err)
	}
	return os.Open(filepath.Join(source, overlayID))
}

func (cs *OverlayContentStore) Remove(ctx context.Context, overlayID string) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	return os.Remove(filepath.Join(source, overlayID))
}

func (*OverlayContentStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_CONTENT_PATH", os.UserConfigDir, "overlay.v4")
	if err != nil {
		return "", fmt.Errorf("search overlay content path: %w", err)
	}
	return path, nil
}

var _ overlay.ContentStore = (*OverlayContentStore)(nil)
