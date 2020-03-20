package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
)

func ExampleList_url() {
	tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	defer os.RemoveAll(tmp)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755)
	_ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755)
	yml := strings.NewReader("githubHost: example.com\nroots:\n  - " + tmp)

	config, err := env.GetAccess(yml, env.EnvarPrefix)
	if err != nil {
		panic(err)
	}
	if err := command.List(&config,
		gogh.URLFormatter(),
		true,
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
	yml := strings.NewReader("githubHost: example.com\nroots:\n  - " + tmp)

	config, err := env.GetAccess(yml, env.EnvarPrefix)
	if err != nil {
		panic(err)
	}
	fmter, err := gogh.CustomFormatter("{{short .}};;{{relative .}}")
	if err != nil {
		panic(err)
	}
	if err := command.List(&config,
		fmter,
		true,
		"",
	); err != nil {
		panic(err)
	}

	// Unordered output:
	// gogh;;example.com/kyoh86/gogh
	// name;;example.com/owner/name
}
