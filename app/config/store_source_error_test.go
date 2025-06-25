package config_test

import (
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
)

func TestStoreSourceErrors(t *testing.T) {
	// テスト前にAppContextPathFuncを保存
	originalFunc := testtarget.AppContextPathFunc
	defer func() {
		testtarget.AppContextPathFunc = originalFunc
	}()

	// AppContextPathFuncをエラーを返すようにモック
	testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return "", errors.New("test error")
	}

	t.Run("HookStore.Source error", func(t *testing.T) {
		store := testtarget.NewHookStore()
		_, err := store.Source()
		if err == nil {
			t.Error("HookStore.Source() error = nil, want error")
		}
	})

	t.Run("OverlayStore.Source error", func(t *testing.T) {
		store := testtarget.NewOverlayStore()
		_, err := store.Source()
		if err == nil {
			t.Error("OverlayStore.Source() error = nil, want error")
		}
	})

	t.Run("ScriptStore.Source error", func(t *testing.T) {
		store := testtarget.NewScriptStore()
		_, err := store.Source()
		if err == nil {
			t.Error("ScriptStore.Source() error = nil, want error")
		}
	})

	t.Run("ExtraStore.Source error", func(t *testing.T) {
		store := testtarget.NewExtraStore()
		_, err := store.Source()
		if err == nil {
			t.Error("ExtraStore.Source() error = nil, want error")
		}
	})

	t.Run("TokenStore.Source error", func(t *testing.T) {
		store := testtarget.NewTokenStore()
		_, err := store.Source()
		if err == nil {
			t.Error("TokenStore.Source() error = nil, want error")
		}
	})

	t.Run("WorkspaceStore.Source error", func(t *testing.T) {
		store := testtarget.NewWorkspaceStore()
		_, err := store.Source()
		if err == nil {
			t.Error("WorkspaceStore.Source() error = nil, want error")
		}
	})

	t.Run("DefaultNameStore.Source error", func(t *testing.T) {
		store := testtarget.NewDefaultNameStore()
		_, err := store.Source()
		if err == nil {
			t.Error("DefaultNameStore.Source() error = nil, want error")
		}
	})

	t.Run("FlagsStore.Source error", func(t *testing.T) {
		store := testtarget.NewFlagsStore()
		_, err := store.Source()
		if err == nil {
			t.Error("FlagsStore.Source() error = nil, want error")
		}
	})

	t.Run("TokenStoreV0.Source error", func(t *testing.T) {
		store := testtarget.NewTokenStoreV0()
		_, err := store.Source()
		if err == nil {
			t.Error("TokenStoreV0.Source() error = nil, want error")
		}
	})

	t.Run("WorkspaceStoreV0.Source error", func(t *testing.T) {
		store := testtarget.NewWorkspaceStoreV0()
		_, err := store.Source()
		if err == nil {
			t.Error("WorkspaceStoreV0.Source() error = nil, want error")
		}
	})

	t.Run("DefaultNameStoreV0.Source error", func(t *testing.T) {
		store := testtarget.NewDefaultNameStoreV0()
		_, err := store.Source()
		if err == nil {
			t.Error("DefaultNameStoreV0.Source() error = nil, want error")
		}
	})

	t.Run("FlagsStoreV0.Source error", func(t *testing.T) {
		store := testtarget.NewFlagsStoreV0()
		_, err := store.Source()
		if err == nil {
			t.Error("FlagsStoreV0.Source() error = nil, want error")
		}
	})

	t.Run("OverlayContentStore.Source error", func(t *testing.T) {
		store := testtarget.NewOverlayContentStore()
		_, err := store.Source()
		if err == nil {
			t.Error("OverlayContentStore.Source() error = nil, want error")
		}
	})

	t.Run("ScriptSourceStore.Source error", func(t *testing.T) {
		store := testtarget.NewScriptSourceStore()
		_, err := store.Source()
		if err == nil {
			t.Error("ScriptSourceStore.Source() error = nil, want error")
		}
	})
}
