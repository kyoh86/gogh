package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/internal/context"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhere(t *testing.T) {
	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)
	defer os.RemoveAll(root2)

	proj1 := filepath.Join(root1, "github.com", "kyoh86", "vim-gogh", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	proj2 := filepath.Join(root2, "github.com", "kyoh86", "gogh", ".git")
	require.NoError(t, os.MkdirAll(proj2, 0755))
	proj3 := filepath.Join(root2, "github.com", "kyoh85", "test", ".git")
	require.NoError(t, os.MkdirAll(proj3, 0755))

	assert.EqualError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, false, "gogh"), "try more precise name")

	assert.EqualError(t, Where(&context.MockContext{
		MStderr:     testutil.DefaultErrorWriter,
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, false, "gogh"), "error writer")

	assert.EqualError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, true, true, "gogh"), "project not found")

	assert.EqualError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, false, "noone"), "project not found")

	assert.NoError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, true, false, "gogh"))

	assert.NoError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, true, "gogh"))

	assert.NoError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, true, "test"))

	assert.NoError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, true, true, "vim-gogh"))

	assert.EqualError(t, Where(&context.MockContext{
		MStdout:     testutil.DefaultErrorWriter,
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, true, true, "vim-gogh"), "error writer")

	assert.EqualError(t, Where(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, true, ".."), "'.' or '..' is reserved name")

	assert.EqualError(t, Where(&context.MockContext{
		MRoot:       []string{"/\x00"},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
	}, false, false, "gogh"), "stat /\x00: invalid argument")

}
