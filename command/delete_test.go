package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/internal/context"
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

	yes := new(yesman)
	assert.NoError(t, command.Delete(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
		MStdin:      yes,
	}, false, "gogh-test-2"), "delete proj2 explicitly")
	{
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj2, 0755))
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	}

	assert.NoError(t, command.Delete(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
		MStdin:      yes,
	}, false, "3"), "delete proj3 with fuzzy")
	{
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj3, 0755))
	}

	assert.NoError(t, command.Delete(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
		MStdin:      yes,
	}, true, "test"), "delete proj1 in primary")
	{
		var err error
		_, err = os.Stat(proj1)
		assert.True(t, os.IsNotExist(err))
		require.NoError(t, os.MkdirAll(proj1, 0755))
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	}

	assert.EqualError(t, command.Delete(&context.MockContext{
		MRoot:       []string{root1, root2},
		MGitHubHost: "github.com",
		MGitHubUser: "kyoh86",
		MStdin:      yes,
	}, true, "foobar"), "any projects did not matched for \"foobar\"")
	{
		var err error
		_, err = os.Stat(proj1)
		assert.NoError(t, err)
		_, err = os.Stat(proj2)
		assert.NoError(t, err)
		_, err = os.Stat(proj3)
		assert.NoError(t, err)
	}
}

type yesman struct {
	odd bool
}

func (y *yesman) Read(b []byte) (int, error) {
	if y.odd {
		b[0] = '\n'
	} else {
		b[0] = 'Y'
	}
	y.odd = !y.odd
	return 1, nil
}
