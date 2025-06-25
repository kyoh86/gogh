package config_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/store"
)

func TestAppContextPath(t *testing.T) {
	testCases := []struct {
		name     string
		envar    string
		envValue string
		getDir   func() (string, error)
		rel      []string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "environment variable set",
			envar:    "TEST_APP_PATH",
			envValue: "/custom/path",
			getDir:   func() (string, error) { return "/config", nil },
			rel:      []string{"test.yaml"},
			wantPath: "/custom/path",
			wantErr:  false,
		},
		{
			name:     "no environment variable, use getDir",
			envar:    "UNSET_VAR",
			envValue: "",
			getDir:   func() (string, error) { return filepath.Join("/home", "user", ".config"), nil },
			rel:      []string{"config.yaml"},
			wantPath: filepath.Join("/home", "user", ".config", "gogh", "config.yaml"),
			wantErr:  false,
		},
		{
			name:     "multiple relative paths",
			envar:    "UNSET_VAR",
			envValue: "",
			getDir:   func() (string, error) { return filepath.Join("/home", "user", ".config"), nil },
			rel:      []string{"subdir", "config.yaml"},
			wantPath: filepath.Join("/home", "user", ".config", "gogh", "subdir", "config.yaml"),
			wantErr:  false,
		},
		{
			name:     "getDir returns error",
			envar:    "UNSET_VAR",
			envValue: "",
			getDir:   func() (string, error) { return "", errors.New("no home dir") },
			rel:      []string{"config.yaml"},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable if needed
			if tc.envValue != "" {
				os.Setenv(tc.envar, tc.envValue)
				defer os.Unsetenv(tc.envar)
			}

			got, err := testtarget.AppContextPath(tc.envar, tc.getDir, tc.rel...)
			if (err != nil) != tc.wantErr {
				t.Errorf("AppContextPath() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.wantPath {
				t.Errorf("AppContextPath() = %v, want %v", got, tc.wantPath)
			}
		})
	}
}

// Mock loader for testing LoadAlternative
type mockLoader[T store.Content] struct {
	shouldFail bool
	failError  error
	value      T
}

func (m mockLoader[T]) Load(ctx context.Context, initial func() T) (T, error) {
	if m.shouldFail {
		return m.value, m.failError
	}
	return m.value, nil
}

func (m mockLoader[T]) Source() (string, error) {
	return "mock", nil
}

// Mock content type for testing
type mockContent struct {
	data string
}

func (m mockContent) HasChanges() bool {
	return false
}

func (m mockContent) MarkSaved() {
	// no-op
}

func TestLoadAlternative(t *testing.T) {
	ctx := context.Background()

	initial := func() mockContent {
		return mockContent{data: "initial"}
	}

	testCases := []struct {
		name     string
		loaders  []store.Loader[mockContent]
		wantData string
		wantErr  bool
	}{
		{
			name: "first loader succeeds",
			loaders: []store.Loader[mockContent]{
				mockLoader[mockContent]{shouldFail: false, value: mockContent{data: "first"}},
				mockLoader[mockContent]{shouldFail: false, value: mockContent{data: "second"}},
			},
			wantData: "first",
			wantErr:  false,
		},
		{
			name: "first loader not exist, second succeeds",
			loaders: []store.Loader[mockContent]{
				mockLoader[mockContent]{shouldFail: true, failError: os.ErrNotExist},
				mockLoader[mockContent]{shouldFail: false, value: mockContent{data: "second"}},
			},
			wantData: "second",
			wantErr:  false,
		},
		{
			name: "first loader fs.ErrNotExist, second succeeds",
			loaders: []store.Loader[mockContent]{
				mockLoader[mockContent]{shouldFail: true, failError: fs.ErrNotExist},
				mockLoader[mockContent]{shouldFail: false, value: mockContent{data: "second"}},
			},
			wantData: "second",
			wantErr:  false,
		},
		{
			name: "loader returns other error",
			loaders: []store.Loader[mockContent]{
				mockLoader[mockContent]{shouldFail: true, failError: errors.New("load error")},
			},
			wantData: "",
			wantErr:  true,
		},
		{
			name: "all loaders fail with not exist, use initial",
			loaders: []store.Loader[mockContent]{
				mockLoader[mockContent]{shouldFail: true, failError: os.ErrNotExist},
				mockLoader[mockContent]{shouldFail: true, failError: fs.ErrNotExist},
			},
			wantData: "initial",
			wantErr:  false,
		},
		{
			name:     "no loaders, use initial",
			loaders:  []store.Loader[mockContent]{},
			wantData: "initial",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := testtarget.LoadAlternative(ctx, initial, tc.loaders...)
			if (err != nil) != tc.wantErr {
				t.Errorf("LoadAlternative() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got.data != tc.wantData {
				t.Errorf("LoadAlternative() = %v, want %v", got.data, tc.wantData)
			}
		})
	}
}
