package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/gogh"
	"github.com/kyoh86/gogh/v4/core/store"
)

// AppContextPath returns the path to the app's configuration file.
//
// If the environment variable `envar` is set, it returns that.
// Specify a function to get the parent directory where the file will be placed, such as os.UserConfigDir.
// The `rel` is the relative path to the file from the dir.
//
// It will make the path that is formed as {getDir()}/{AppName=gogh}/{rel...}`
func AppContextPath(envar string, getDir func() (string, error), rel ...string) (string, error) {
	if env := os.Getenv(envar); env != "" {
		return env, nil
	}
	dir, err := getDir()
	if err != nil {
		return "", fmt.Errorf("search app file dir for %s: %w", rel, err)
	}
	return filepath.Join(append([]string{dir, gogh.AppName}, rel...)...), nil
}

var AppContextPathFunc = AppContextPath

func appHomePath(envar, homeEnvar string, getDir func() (string, error), rel ...string) (string, error) {
	if env := os.Getenv(envar); env != "" {
		return env, nil
	}
	if home := os.Getenv(homeEnvar); home != "" {
		return filepath.Join(append([]string{home}, rel...)...), nil
	}
	return AppContextPathFunc(envar, getDir, rel...)
}

// AppConfigPath returns the path to the app's configuration file.
//
// If envar is set, it returns that. Otherwise, if GOGH_CONFIG_HOME is set, it
// returns the relative path under that directory. Finally, it falls back to
// os.UserConfigDir()/gogh.
func AppConfigPath(envar string, rel ...string) (string, error) {
	return appHomePath(envar, "GOGH_CONFIG_HOME", os.UserConfigDir, rel...)
}

var AppConfigPathFunc = AppConfigPath

// AppCachePath returns the path to the app's cache file.
//
// If envar is set, it returns that. Otherwise, if GOGH_CACHE_HOME is set, it
// returns the relative path under that directory. Finally, it falls back to
// os.UserCacheDir()/gogh.
func AppCachePath(envar string, rel ...string) (string, error) {
	return appHomePath(envar, "GOGH_CACHE_HOME", os.UserCacheDir, rel...)
}

var AppCachePathFunc = AppCachePath

// LoadAlternative loads a value of type T using the provided loaders.
// It tries each loader in order until one succeeds or all fail.
// If all loaders fail, it returns the initial value provided by the initial function.
func LoadAlternative[T store.Content](ctx context.Context, initial func() T, loaders ...store.Loader[T]) (T, error) {
	for i, loader := range loaders {
		svc, err := loader.Load(ctx, initial)
		if os.IsNotExist(err) {
			continue
		}
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			var empty T
			return empty, fmt.Errorf("loading at %dth loader: %w", i+1, err)
		}
		return svc, nil
	}
	return initial(), nil
}
