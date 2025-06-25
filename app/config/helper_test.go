package config_test

import (
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
)

func TestOverlayDir(t *testing.T) {
	// テスト前にAppContextPathFuncを保存
	originalFunc := testtarget.AppContextPathFunc
	defer func() {
		testtarget.AppContextPathFunc = originalFunc
	}()

	t.Run("success", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "/test/path/overlay.v4.toml", nil
		}

		path, err := testtarget.OverlayDir()
		if err != nil {
			t.Errorf("OverlayDir() error = %v, want nil", err)
		}
		if path != "/test/path/overlay.v4.toml" {
			t.Errorf("OverlayDir() = %v, want /test/path/overlay.v4.toml", path)
		}
	})

	t.Run("error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("test error")
		}

		_, err := testtarget.OverlayDir()
		if err == nil {
			t.Error("OverlayDir() error = nil, want error")
		}
	})
}

func TestScriptDirError(t *testing.T) {
	// テスト前にAppContextPathFuncを保存
	originalFunc := testtarget.AppContextPathFunc
	defer func() {
		testtarget.AppContextPathFunc = originalFunc
	}()

	t.Run("success", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "/test/path/script.v4.toml", nil
		}

		path, err := testtarget.ScriptDir()
		if err != nil {
			t.Errorf("ScriptDir() error = %v, want nil", err)
		}
		if path != "/test/path/script.v4.toml" {
			t.Errorf("ScriptDir() = %v, want /test/path/script.v4.toml", path)
		}
	})

	t.Run("error", func(t *testing.T) {
		testtarget.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", errors.New("test error")
		}

		_, err := testtarget.ScriptDir()
		if err == nil {
			t.Error("ScriptDir() error = nil, want error")
		}
	})
}
