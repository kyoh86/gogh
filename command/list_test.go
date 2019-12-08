package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleList_url() {
	tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	defer os.RemoveAll(tmp)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755)

	if err := command.List(&config.Config{
		VRoot:  []string{tmp},
		GitHub: config.GitHubConfig{Host: "example.com"},
	},
		gogh.URLFormatter(),
		true,
		false,
		"",
	); err != nil {
		panic(err)
	}
	// Unordered output:
	// https://example.com/kyoh86/gogh
	// https://example.com/owner/name
}

func ExampleList_custom() {
	tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	defer os.RemoveAll(tmp)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755)

	fmter, err := gogh.CustomFormatter("{{short .}};;{{relative .}}")
	if err != nil {
		panic(err)
	}
	if err := command.List(&config.Config{
		VRoot:  []string{tmp},
		GitHub: config.GitHubConfig{Host: "example.com"},
	},
		fmter,
		true,
		false,
		"",
	); err != nil {
		panic(err)
	}
	// Unordered output:
	// gogh;;example.com/kyoh86/gogh
	// name;;example.com/owner/name
}

func TestList(t *testing.T) {
	tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	defer os.RemoveAll(tmp)
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755))

	assert.Error(t, command.List(&config.Config{
		VRoot:  []string{tmp, "/\x00"},
		GitHub: config.GitHubConfig{Host: "example.com"},
	},
		gogh.URLFormatter(),
		false,
		false,
		"",
	), "invalid root")
}
