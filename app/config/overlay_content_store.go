package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// overlayKey generates a unique hash string from Overlay struct content.
func overlayKey(ov overlay.Overlay) (string, error) {
	b, err := json.Marshal(ov)
	if err != nil {
		return "", fmt.Errorf("failed to marshal overlay: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil // 64 chars
}

func (cs *OverlayContentStore) SaveContent(ctx context.Context, ov overlay.Overlay, content io.Reader) (string, error) {
	source, err := cs.Source()
	if err != nil {
		return "", fmt.Errorf("get content source: %w", err)
	}
	if err := os.MkdirAll(source, 0755); err != nil {
		return "", fmt.Errorf("failed to create content directory: %w", err)
	}
	fileName, err := overlayKey(ov)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(source, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create content file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, content); err != nil {
		return "", fmt.Errorf("failed to write content: %w", err)
	}
	return fileName, nil
}

func (cs *OverlayContentStore) OpenContent(ctx context.Context, location string) (io.ReadCloser, error) {
	source, err := cs.Source()
	if err != nil {
		return nil, fmt.Errorf("get content source: %w", err)
	}
	return os.Open(filepath.Join(source, location))
}

func (cs *OverlayContentStore) RemoveContent(ctx context.Context, location string) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	return os.Remove(filepath.Join(source, location))
}

func (*OverlayContentStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_OVERLAY_CONTENT_PATH", os.UserConfigDir, "overlay.v4")
	if err != nil {
		return "", fmt.Errorf("search overlay content path: %w", err)
	}
	return path, nil
}

var _ overlay.ContentStore = (*OverlayContentStore)(nil)
