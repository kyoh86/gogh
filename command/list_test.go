package command_test

import (
	"testing"
)

func ExampleList_url() {
	// UNDONE: mock and example
	// tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	// defer os.RemoveAll(tmp)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755)

	// if err := command.List(&config.Config{
	// 	VRoot:  []string{tmp},
	// 	GitHub: config.GitHubConfig{Host: "example.com"},
	// },
	// 	gogh.URLFormatter(),
	// 	true,
	// 	"",
	// ); err != nil {
	// 	panic(err)
	// }

	// Unordered output:
	// https://example.com/kyoh86/gogh
	// https://example.com/owner/name
}

func ExampleList_custom() {
	// UNDONE: mock and example
	// tmp, _ := ioutil.TempDir(os.TempDir(), "gogh-test")
	// defer os.RemoveAll(tmp)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "kyoh86", "gogh", ".git"), 0755)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "name", ".git"), 0755)
	// _ = os.MkdirAll(filepath.Join(tmp, "example.com", "owner", "empty"), 0755)
	//
	// fmter, err := gogh.CustomFormatter("{{short .}};;{{relative .}}")
	// if err != nil {
	// 	panic(err)
	// }
	// if err := command.List(&config.Config{
	// 	VRoot:  []string{tmp},
	// 	GitHub: config.GitHubConfig{Host: "example.com"},
	// },
	// 	fmter,
	// 	true,
	// 	"",
	// ); err != nil {
	// 	panic(err)
	// }

	// Unordered output:
	// gogh;;example.com/kyoh86/gogh
	// name;;example.com/owner/name
}

func TestList(t *testing.T) {
	/* UNDONE:
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
		"",
	), "invalid root")
	*/
}
