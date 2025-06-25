package config_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestOverlayContentStoreSaveErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("Save with source error", func(t *testing.T) {
		// テスト前にAppContextPathFuncを保存
		originalFunc := testtarget.AppContextPathFunc
		defer func() {
			testtarget.AppContextPathFunc = originalFunc
		}()

		// AppContextPathFuncをエラーを返すようにモック
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("source error")
		}

		store := testtarget.NewOverlayContentStore()
		err := store.Save(ctx, "test-id", strings.NewReader("test content"))
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})

	t.Run("Open with source error", func(t *testing.T) {
		// テスト前にAppContextPathFuncを保存
		originalFunc := testtarget.AppContextPathFunc
		defer func() {
			testtarget.AppContextPathFunc = originalFunc
		}()

		// AppContextPathFuncをエラーを返すようにモック
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("source error")
		}

		store := testtarget.NewOverlayContentStore()
		_, err := store.Open(ctx, "test-id")
		if err == nil {
			t.Error("Open() error = nil, want error")
		}
	})

	t.Run("Remove with source error", func(t *testing.T) {
		// テスト前にAppContextPathFuncを保存
		originalFunc := testtarget.AppContextPathFunc
		defer func() {
			testtarget.AppContextPathFunc = originalFunc
		}()

		// AppContextPathFuncをエラーを返すようにモック
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("source error")
		}

		store := testtarget.NewOverlayContentStore()
		err := store.Remove(ctx, "test-id")
		if err == nil {
			t.Error("Remove() error = nil, want error")
		}
	})
}

func TestScriptSourceStoreSaveErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("Save with source error", func(t *testing.T) {
		// テスト前にAppContextPathFuncを保存
		originalFunc := testtarget.AppContextPathFunc
		defer func() {
			testtarget.AppContextPathFunc = originalFunc
		}()

		// AppContextPathFuncをエラーを返すようにモック
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("source error")
		}

		store := testtarget.NewScriptSourceStore()
		err := store.Save(ctx, "test-id", strings.NewReader("test script"))
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})
}

func TestStoreSaveListErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	// Create a temporary directory
	tempDir := t.TempDir()

	t.Run("ExtraStore.Save with List error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return filepath.Join(tempDir, "extra.toml"), nil
		}

		store := testtarget.NewExtraStore()
		mockService := extra_mock.NewMockExtraService(ctrl)

		mockService.EXPECT().HasChanges().Return(true)
		mockService.EXPECT().List(gomock.Any()).Return(func(yield func(*extra.Extra, error) bool) {
			yield(nil, errors.New("list error"))
		})

		err := store.Save(ctx, mockService, false)
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})

	t.Run("HookStore.Save with List yield error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return filepath.Join(tempDir, "hook.toml"), nil
		}

		store := testtarget.NewHookStore()
		mockService := hook_mock.NewMockHookService(ctrl)

		mockService.EXPECT().HasChanges().Return(true)
		mockService.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
			yield(nil, errors.New("yield error"))
		})

		err := store.Save(ctx, mockService, false)
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})

	t.Run("OverlayStore.Save with List yield error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return filepath.Join(tempDir, "overlay.toml"), nil
		}

		store := testtarget.NewOverlayStore()
		mockService := overlay_mock.NewMockOverlayService(ctrl)

		mockService.EXPECT().HasChanges().Return(true)
		mockService.EXPECT().List().Return(func(yield func(overlay.Overlay, error) bool) {
			yield(nil, errors.New("yield error"))
		})

		err := store.Save(ctx, mockService, false)
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})

	t.Run("ScriptStore.Save with List yield error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return filepath.Join(tempDir, "script.toml"), nil
		}

		store := testtarget.NewScriptStore()
		mockService := script_mock.NewMockScriptService(ctrl)

		mockService.EXPECT().HasChanges().Return(true)
		mockService.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
			yield(nil, errors.New("yield error"))
		})

		err := store.Save(ctx, mockService, false)
		if err == nil {
			t.Error("Save() error = nil, want error")
		}
	})
}

// mockReader simulates an io.Reader that always returns an error
type mockErrorReader struct{}

func (m *mockErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestOverlayContentStoreSaveIOError(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Set up path to use temp directory
	testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return filepath.Join(tempDir, "overlay.v4"), nil
	}

	store := testtarget.NewOverlayContentStore()

	// Test with a reader that returns an error
	err := store.Save(ctx, "test-id", &mockErrorReader{})
	if err == nil {
		t.Error("Save() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "failed to write content") {
		t.Errorf("Save() error = %v, want error containing 'failed to write content'", err)
	}
}

func TestScriptSourceStoreSaveIOError(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Set up path to use temp directory
	testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return filepath.Join(tempDir, "script.v4"), nil
	}

	store := testtarget.NewScriptSourceStore()

	// Test with a reader that returns an error
	err := store.Save(ctx, "test-id", &mockErrorReader{})
	if err == nil {
		t.Error("Save() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "failed to write script") {
		t.Errorf("Save() error = %v, want error containing 'failed to write script'", err)
	}
}
