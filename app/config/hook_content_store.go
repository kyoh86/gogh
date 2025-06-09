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

	"github.com/kyoh86/gogh/v4/core/hook"
)

type HookContentStore struct{}

func NewHookContentStore() *HookContentStore { return &HookContentStore{} }

func hookKey(h hook.Hook) (string, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("failed to marshal hook: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func (cs *HookContentStore) SaveScript(ctx context.Context, h hook.Hook, content io.Reader) (string, error) {
	source, err := cs.Source()
	if err != nil {
		return "", fmt.Errorf("get content source: %w", err)
	}
	if err := os.MkdirAll(source, 0755); err != nil {
		return "", fmt.Errorf("failed to create content directory: %w", err)
	}
	fileName, err := hookKey(h)
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

func (cs *HookContentStore) OpenScript(ctx context.Context, location string) (io.ReadCloser, error) {
	source, err := cs.Source()
	if err != nil {
		return nil, fmt.Errorf("get content source: %w", err)
	}
	return os.Open(filepath.Join(source, location))
}

func (cs *HookContentStore) RemoveScript(ctx context.Context, location string) error {
	source, err := cs.Source()
	if err != nil {
		return fmt.Errorf("get content source: %w", err)
	}
	return os.Remove(filepath.Join(source, location))
}

func (*HookContentStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_HOOK_CONTENT_PATH", os.UserConfigDir, "hook.v4")
	if err != nil {
		return "", fmt.Errorf("search hook content path: %w", err)
	}
	return path, nil
}

var _ hook.HookScriptStore = (*HookContentStore)(nil)
