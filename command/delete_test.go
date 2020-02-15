package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/kyoh86/gogh/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)
	defer os.RemoveAll(root1)

	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)
	defer os.RemoveAll(root2)

	proj1 := filepath.Join(root1, "github.com", "kyoh86", "gogh-test-1", ".git")
	require.NoError(t, os.MkdirAll(proj1, 0755))
	proj2 := filepath.Join(root2, "github.com", "kyoh86", "gogh-test-2", ".git")
	require.NoError(t, os.MkdirAll(proj2, 0755))
	proj3 := filepath.Join(root2, "github.com", "kyoh85", "gogh-test-3", ".git")
	require.NoError(t, os.MkdirAll(proj3, 0755))

	t.Run("delete proj2 explicitly", func(t *testing.T) {
		teardown := testutil.Stubin(t, []byte("y\n"))
		defer teardown()

		assert.NoError(t, command.Delete(&context.MockContext{
			MRoot:       []string{root1, root2},
			MGitHubHost: "github.com",
			MGitHubUser: "kyoh86",
		}, false, "gogh-test-2"))
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj2, 0755))
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})

	t.Run("delete proj3 with fuzzy", func(t *testing.T) {
		teardown := testutil.Stubin(t, []byte("y\n"))
		defer teardown()

		assert.NoError(t, command.Delete(&context.MockContext{
			MRoot:       []string{root1, root2},
			MGitHubHost: "github.com",
			MGitHubUser: "kyoh86",
		}, false, "3"))
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj3, 0755))
	})

	t.Run("delete proj1 in primary", func(t *testing.T) {
		teardown := testutil.Stubin(t, []byte("y\n"))
		defer teardown()

		assert.NoError(t, command.Delete(&context.MockContext{
			MRoot:       []string{root1, root2},
			MGitHubHost: "github.com",
			MGitHubUser: "kyoh86",
		}, true, "test"))
		var err error
		_, err = os.Stat(proj1)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj1, 0755))
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})

	t.Run("did not match", func(t *testing.T) {
		assert.EqualError(t, command.Delete(&context.MockContext{
			MRoot:       []string{root1, root2},
			MGitHubHost: "github.com",
			MGitHubUser: "kyoh86",
		}, true, "foobar"), "any projects did not matched for \"foobar\"")
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	})
}
